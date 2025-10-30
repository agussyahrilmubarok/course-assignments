package com.example.timesale.service.v1;

import com.example.timesale.aop.TimeSaleMetered;
import com.example.timesale.domain.Product;
import com.example.timesale.domain.TimeSale;
import com.example.timesale.domain.TimeSaleOrder;
import com.example.timesale.model.TimeSaleDTO;
import com.example.timesale.repos.ProductRepository;
import com.example.timesale.repos.TimeSaleOrderRepository;
import com.example.timesale.repos.TimeSaleRepository;
import com.example.timesale.service.TimeSaleService;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.time.LocalDateTime;
import java.util.UUID;

@Service("TimeSaleServiceImplV1")
@Slf4j
@RequiredArgsConstructor
public class TimeSaleServiceImplV1 implements TimeSaleService {

    private final TimeSaleRepository timeSaleRepository;
    private final ProductRepository productRepository;
    private final TimeSaleOrderRepository timeSaleOrderRepository;

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
        log.info("TimeSale created successfully with ID: {}", saved.getId());
        return saved;
    }

    @Override
    @Transactional(readOnly = true)
    public TimeSale findById(String timeSaleId) {
        return timeSaleRepository.findById(timeSaleId)
                .orElseThrow(() -> {
                    log.warn("TimeSale with ID {} not found", timeSaleId);
                    return new IllegalArgumentException("Time sale not found");
                });
    }

    @Transactional(readOnly = true)
    public Page<TimeSale> findAllOngoing(Pageable pageable) {
        LocalDateTime now = LocalDateTime.now();
        return timeSaleRepository.findAllByStartAtBeforeAndEndAtAfterAndStatus(now, TimeSale.Status.ACTIVE, pageable);
    }

    @Override
    @Transactional
    @TimeSaleMetered(version = "v1")
    public TimeSale purchase(TimeSaleDTO.PurchaseRequest request, String userId) {
        TimeSale timeSale = timeSaleRepository.findByIdWithPessimisticLock(request.getTimeSaleId())
                .orElseThrow(() -> new IllegalArgumentException("TimeSale not found"));
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

        return timeSale;
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
        if (quantity <= 0) {
            throw new IllegalArgumentException("Quantity must be positive");
        }
        if (discountPrice <= 0) {
            throw new IllegalArgumentException("Discount price must be positive");
        }
    }
}
