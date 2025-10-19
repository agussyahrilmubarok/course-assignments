package com.example.order.service;

import com.example.order.domain.ProductOrder;
import com.example.order.exception.InsufficientStockException;
import com.example.order.model.OrderDTO;
import com.example.order.model.ProductDTO;
import com.example.order.repos.ProductOrderRepository;
import jakarta.persistence.EntityNotFoundException;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Service;

import java.util.List;
import java.util.UUID;
import java.util.stream.Collectors;

@Service
@Slf4j
@RequiredArgsConstructor
public class OrderServiceImpl implements OrderService {

    private final ProductOrderRepository productOrderRepository;
    private final CatalogClient catalogClient;

    @Override
    public OrderDTO.StartOrderResponse startOrder(OrderDTO.StartOrderRequest payload) {
        ProductDTO.Response product = catalogClient.getProductById(payload.getProductId());
        if (product.getStockCount() < payload.getCount()) {
            throw new InsufficientStockException("Product is out of stock or quantity requested");
        }

        ProductOrder order = new ProductOrder();
        order.setId(UUID.randomUUID().toString());
        order.setUserId(payload.getUserId());
        order.setProductId(product.getId());
        order.setCount(payload.getCount());
        order.setPaymentId("");
        order.setOrderStatus(ProductOrder.Status.CREATED);
        order = productOrderRepository.save(order);

        return OrderDTO.StartOrderResponse.builder()
                .orderId(order.getId())
                .paymentUrl("")
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
