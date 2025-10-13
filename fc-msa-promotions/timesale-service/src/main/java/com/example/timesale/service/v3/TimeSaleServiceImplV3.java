package com.example.timesale.service.v3;

import com.example.timesale.domain.Product;
import com.example.timesale.domain.TimeSale;
import com.example.timesale.domain.TimeSaleOrder;
import com.example.timesale.model.TimeSaleDTO;
import com.example.timesale.repos.ProductRepository;
import com.example.timesale.repos.TimeSaleOrderRepository;
import com.example.timesale.repos.TimeSaleRepository;
import com.example.timesale.service.TimeSaleAsyncService;
import com.example.timesale.service.TimeSaleService;
import com.example.timesale.service.v3.component.KafkaProducer;
import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.redisson.api.RAtomicLong;
import org.redisson.api.RBucket;
import org.redisson.api.RLock;
import org.redisson.api.RedissonClient;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.time.LocalDateTime;
import java.util.UUID;
import java.util.concurrent.TimeUnit;

@Service("TimeSaleServiceImplV3")
@Slf4j
@RequiredArgsConstructor
public class TimeSaleServiceImplV3 implements TimeSaleService, TimeSaleAsyncService {

    private static final String TIME_SALE_KEY = "time-sale:";
    private static final String TIME_SALE_LOCK = "time-sale-lock:";
    private static final long WAIT_TIME = 3L;
    private static final long LEASE_TIME = 3L;
    private static final String QUEUE_KEY = "time-sale-queue:";
    private static final String TOTAL_REQUESTS_KEY = "time-sale-total-requests:";
    private static final String RESULT_PREFIX = "purchase-result:";

    private final TimeSaleRepository timeSaleRepository;
    private final ProductRepository productRepository;
    private final TimeSaleOrderRepository timeSaleOrderRepository;
    private final RedissonClient redissonClient;
    private final ObjectMapper objectMapper;
    private final KafkaProducer kafkaProducer;

    @Override
    @Transactional
    public TimeSale create(TimeSaleDTO.CreateRequest request) {
        this.validateTimeSale(request.getQuantity(), request.getDiscountPrice(),
                request.getStartAt(), request.getEndAt());

        Product product = this.saveProduct(request.getProduct());

        TimeSale timeSale = new TimeSale();
        timeSale.setId(UUID.randomUUID().toString());
        timeSale.setQuantity(request.getQuantity());
        timeSale.setRemainingQuantity(request.getQuantity());
        timeSale.setDiscountPrice(request.getDiscountPrice());
        timeSale.setStartAt(request.getStartAt());
        timeSale.setEndAt(request.getEndAt());
        timeSale.setStatus(TimeSale.Status.ACTIVE);
        timeSale.setProduct(product);

        TimeSale saved = timeSaleRepository.save(timeSale);
        this.setTimeSaleCache(saved);
        log.info("TimeSale created successfully with ID: {}", saved.getId());
        return saved;
    }

    @Override
    @Transactional(readOnly = true)
    public TimeSale findById(String timeSaleId) {
        return this.getTimeSaleCache(timeSaleId);
    }

    @Transactional(readOnly = true)
    public Page<TimeSale> findAllOngoing(Pageable pageable) {
        LocalDateTime now = LocalDateTime.now();
        return timeSaleRepository.findAllByStartAtBeforeAndEndAtAfterAndStatus(now, TimeSale.Status.ACTIVE, pageable);
    }

    @Override
    @Transactional
    public TimeSale purchase(TimeSaleDTO.PurchaseRequest request, String userId) {
        RLock lock = redissonClient.getLock(TIME_SALE_LOCK + request.getTimeSaleId());
        if (lock == null) {
            throw new IllegalStateException("Failed to acquire lock instance.");
        }

        boolean isLocked = false;

        try {
            isLocked = lock.tryLock(WAIT_TIME, LEASE_TIME, TimeUnit.SECONDS);
            if (!isLocked) {
                throw new IllegalStateException("Could not acquire lock for time sale purchase.");
            }

            TimeSale timeSale = this.getTimeSaleCache(request.getTimeSaleId());
            timeSale.purchase(request.getQuantity());
            timeSale = timeSaleRepository.save(timeSale);
            log.info("Purchase successful. Remaining quantity: {}", timeSale.getRemainingQuantity());

            TimeSaleOrder timeSaleOrder = new TimeSaleOrder();
            timeSaleOrder.setId(UUID.randomUUID().toString());
            timeSaleOrder.setUserId(userId);
            timeSaleOrder.setQuantity(request.getQuantity());
            timeSaleOrder.setDiscountPrice(timeSale.getDiscountPrice());
            timeSaleOrder.setStatus(TimeSaleOrder.Status.PENDING);
            timeSaleOrder.setTimeSale(timeSale);
            timeSaleOrderRepository.save(timeSaleOrder);
            log.debug("TimeSaleOrder saved with ID: {}", timeSaleOrder.getId());

            this.setTimeSaleCache(timeSale);

            return timeSale;
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
            throw new IllegalStateException("Thread was interrupted while waiting for lock.");
        } catch (Exception e) {
            log.error("Error during purchase", e);
            throw new RuntimeException("Failed to process purchase", e);
        } finally {
            if (isLocked) {
                try {
                    lock.unlock();
                } catch (Exception e) {
                    log.error("Failed to unlock", e);
                }
            }
        }
    }

    /**
     * - Handles purchase requests asynchronously via Kafka
     * - Manages a request queue using Redis
     * - Uses Redisson for concurrency control in distributed environments
     */
    @Override
    public String purchaseRequest(TimeSaleDTO.PurchaseRequest request, String userId) {
        // Generate a unique request ID
        String requestId = UUID.randomUUID().toString();
        // Create the purchase request message
        TimeSaleDTO.PurchaseRequestMessage message = TimeSaleDTO.PurchaseRequestMessage.builder()
                .requestId(requestId)
                .timeSaleId(request.getTimeSaleId())
                .userId(userId)
                .quantity(request.getQuantity())
                .build();

        // Store initial status in Redis
        RBucket<String> resultBucket = redissonClient.getBucket(RESULT_PREFIX + requestId);
        resultBucket.set("PENDING");
        // Add to queue and increment total request count
        String queueKey = QUEUE_KEY + request.getTimeSaleId();
        String totalKey = TOTAL_REQUESTS_KEY + request.getTimeSaleId();
        RBucket<String> queueBucket = redissonClient.getBucket(queueKey);
        queueBucket.set(requestId);
        RAtomicLong totalCounter = redissonClient.getAtomicLong(totalKey);
        totalCounter.incrementAndGet();

        // Send message to Kafka
        kafkaProducer.sendPurchaseRequest(requestId, message);
        return requestId;
    }

    @Override
    public TimeSaleDTO.AsyncPurchaseResponse findPurchaseResult(String timeSaleId, String requestId) {
        RBucket<String> resultBucket = redissonClient.getBucket(RESULT_PREFIX + requestId);
        String result = resultBucket.get();
        String status = result != null ? result : "PENDING";

        // Get queue position and total waiting users if still pending
        Integer queuePosition = null;
        Long totalWaiting = 0L;

        if ("PENDING".equals(status)) {
            queuePosition = this.getQueuePosition(timeSaleId, requestId);
            totalWaiting = this.getTotalWaiting(timeSaleId);
        }

        return TimeSaleDTO.AsyncPurchaseResponse.builder()
                .requestId(requestId)
                .status(status)
                .queuePosition(queuePosition)
                .totalWaiting(totalWaiting)
                .build();
    }

    /**
     * Save the processing result of a purchase request to Redis.
     *
     * @param requestId the request ID
     * @param result    the result status (SUCCESS/FAIL)
     */
    @Override
    public void savePurchaseResult(String requestId, String result) {
        RBucket<String> resultBucket = redissonClient.getBucket(RESULT_PREFIX + requestId);
        resultBucket.set(result);
    }

    /**
     * Remove the processed request from the queue.
     * 1. Remove request ID from the queue
     * 2. Decrease the total pending request count
     *
     * @param timeSaleId the time sale ID
     * @param requestId  the request ID
     */
    @Override
    public void removePurchaseResultFromQueue(String timeSaleId, String requestId) {
        try {
            // Remove request from the queue
            String queueKey = QUEUE_KEY + timeSaleId;
            RBucket<String> queueBucket = redissonClient.getBucket(queueKey);
            String queueValue = queueBucket.get();

            if (queueValue != null && !queueValue.isEmpty()) {
                // Remove the specific request ID from the comma-separated string
                String[] queueValues = queueValue.split(",");
                StringBuilder newQueue = new StringBuilder();
                for (String value : queueValues) {
                    if (!requestId.equals(value)) {
                        if (!newQueue.isEmpty()) {
                            newQueue.append(",");
                        }
                        newQueue.append(value);
                    }
                }
                queueBucket.set(newQueue.toString());
            }

            // Decrease the total waiting count
            String totalKey = TOTAL_REQUESTS_KEY + timeSaleId;
            RAtomicLong totalCounter = redissonClient.getAtomicLong(totalKey);
            totalCounter.decrementAndGet();
        } catch (Exception e) {
            log.error("Failed to remove request from queue: timeSaleId={}, requestId={}", timeSaleId, requestId, e);
        }
    }

    private Product saveProduct(TimeSaleDTO.ProductRequest request) {
        Product product = new Product();
        product.setId(UUID.randomUUID().toString());
        product.setName(request.getName());
        product.setPrice(request.getPrice());
        return productRepository.save(product);
    }

    private void validateTimeSale(Long quantity, Long discountPrice, LocalDateTime startAt, LocalDateTime endAt) {
        if (startAt.isAfter(endAt)) {
            throw new IllegalArgumentException("Start time must be before end time");
        }
        if (quantity == null || quantity <= 0) {
            throw new IllegalArgumentException("Quantity must be positive");
        }
        if (discountPrice == null || discountPrice <= 0) {
            throw new IllegalArgumentException("Discount price must be positive");
        }
    }

    public void setTimeSaleCache(TimeSale timeSale) {
        try {
            String json = objectMapper.writeValueAsString(timeSale);
            RBucket<String> bucket = redissonClient.getBucket(TIME_SALE_KEY + timeSale.getId());
            bucket.set(json);
        } catch (JsonProcessingException e) {
            log.warn("Failed to serialize TimeSale: {}", e.getMessage());
        } catch (Exception e) {
            log.error("Unexpected error while setting cache", e);
        }
    }

    private TimeSale getTimeSaleCache(String timeSaleId) {
        RBucket<String> bucket = redissonClient.getBucket(TIME_SALE_KEY + timeSaleId);
        String json = bucket.get();

        if (json != null) {
            try {
                return objectMapper.readValue(json, TimeSale.class);
            } catch (JsonProcessingException e) {
                log.warn("Failed to deserialize TimeSale from cache: {}", e.getMessage());
            }
        }

        TimeSale timeSale = timeSaleRepository.findById(timeSaleId)
                .orElseThrow(() -> new IllegalArgumentException("TimeSale not found"));
        setTimeSaleCache(timeSale);
        return timeSale;
    }

    /**
     * Retrieves the position of a request in the waiting queue.
     */
    private Integer getQueuePosition(String timeSaleId, String requestId) {
        String queueKey = QUEUE_KEY + timeSaleId;
        RBucket<String> queueBucket = redissonClient.getBucket(queueKey);
        String queueValue = queueBucket.get();

        if (queueValue == null || queueValue.isEmpty()) {
            return null;
        }

        String[] queueValues = queueValue.split(",");
        for (int i = 0; i < queueValues.length; i++) {
            if (requestId.equals(queueValues[i])) {
                return i + 1;
            }
        }
        return null;
    }

    /**
     * Retrieves the total number of pending requests.
     */
    private Long getTotalWaiting(String timeSaleId) {
        String totalKey = TOTAL_REQUESTS_KEY + timeSaleId;
        RAtomicLong totalCounter = redissonClient.getAtomicLong(totalKey);
        return totalCounter.get();
    }
}
