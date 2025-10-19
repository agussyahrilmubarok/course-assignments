package com.example.timesale.service.v2;

import com.example.timesale.domain.Product;
import com.example.timesale.domain.TimeSale;
import com.example.timesale.domain.TimeSaleOrder;
import com.example.timesale.model.TimeSaleDTO;
import com.example.timesale.repos.ProductRepository;
import com.example.timesale.repos.TimeSaleOrderRepository;
import com.example.timesale.repos.TimeSaleRepository;
import com.example.timesale.service.TimeSaleService;
import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
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

@Service("TimeSaleServiceImplV2")
@Slf4j
@RequiredArgsConstructor
public class TimeSaleServiceImplV2 implements TimeSaleService {

    private static final String TIME_SALE_KEY = "time-sale:";
    private static final String TIME_SALE_LOCK = "time-sale-lock:";
    private static final long WAIT_TIME = 3L;
    private static final long LEASE_TIME = 3L;

    private final TimeSaleRepository timeSaleRepository;
    private final ProductRepository productRepository;
    private final TimeSaleOrderRepository timeSaleOrderRepository;
    private final RedissonClient redissonClient;
    private final ObjectMapper objectMapper;

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

    private void setTimeSaleCache(TimeSale timeSale) {
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
}
