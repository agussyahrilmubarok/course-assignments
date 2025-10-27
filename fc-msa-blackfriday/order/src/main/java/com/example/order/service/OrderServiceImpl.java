package com.example.order.service;

import com.example.order.domain.ProductOrder;
import com.example.order.exception.InsufficientStockException;
import com.example.order.exception.PaymentFailedException;
import com.example.order.model.OrderDTO;
import com.example.order.model.PaymentDTO;
import com.example.order.model.ProductDTO;
import com.example.order.repos.ProductOrderRepository;
import jakarta.persistence.EntityNotFoundException;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Service;

import java.math.BigDecimal;
import java.util.List;
import java.util.UUID;
import java.util.stream.Collectors;

@Service
@Slf4j
@RequiredArgsConstructor
public class OrderServiceImpl implements OrderService {

    private final ProductOrderRepository productOrderRepository;
    private final CatalogClient catalogClient;
    private final PaymentClient paymentClient;

    @Override
    public OrderDTO.StartOrderResponse startOrder(OrderDTO.StartOrderRequest payload) {
        ProductDTO.Response product = catalogClient.getProductById(payload.getProductId());
        if (product.getStockCount() < payload.getCount()) {
            throw new InsufficientStockException("Insufficient stock for product ID: " + product.getId());
        }

        BigDecimal totalAmount = BigDecimal.valueOf(product.getPrice()).multiply(BigDecimal.valueOf(payload.getCount()));

        ProductOrder order = new ProductOrder();
        order.setId(UUID.randomUUID().toString());
        order.setUserId(payload.getUserId());
        order.setProductId(product.getId());
        order.setCount(payload.getCount());
        order.setAmount(totalAmount);
        order.setPaymentId(null);
        order.setOrderStatus(ProductOrder.Status.CREATED);
        productOrderRepository.save(order);

        PaymentDTO.Response transaction;
        try {
            transaction = paymentClient.createPayment(
                    PaymentDTO.CreateTransactionRequest.builder()
                            .orderId(order.getId())
                            .amount(totalAmount)
                            .build()
            );
        } catch (Exception ex) {
            log.error("Failed to create payment for order {}: {}", order.getId(), ex.getMessage(), ex);
            throw new PaymentFailedException("Failed to initiate payment for order");
        }

        if (transaction.getId() != null) {
            order.setPaymentId(transaction.getId());
            productOrderRepository.save(order);
        }

        log.info("Order {} created successfully for user {}", order.getId(), payload.getUserId());

        return OrderDTO.StartOrderResponse.builder()
                .orderId(order.getId())
                .paymentUrl(transaction.getPaymentUrl())
                .build();
    }

    @Override
    public OrderDTO.FinishOrderRequest finishOrder(OrderDTO.FinishOrderRequest payload) {
        return null;
    }

    @Override
    public OrderDTO.Response findById(String orderId) {
        ProductOrder order = productOrderRepository.findById(orderId)
                .orElseThrow(() -> {
                    log.warn("Order not found with ID: {}", orderId);
                    return new EntityNotFoundException("Order not found with ID: " + orderId);
                });
        return OrderDTO.Response.from(order);
    }

    @Override
    public List<OrderDTO.Response> findAllByUser(String userId) {
        List<ProductOrder> orders = productOrderRepository.findAllByUserId(userId);
        if (orders.isEmpty()) {
            log.warn("No orders found for user ID: {}", userId);
            throw new EntityNotFoundException("No orders found for user ID: " + userId);
        }
        return orders.stream()
                .map(OrderDTO.Response::from)
                .collect(Collectors.toList());
    }
}
