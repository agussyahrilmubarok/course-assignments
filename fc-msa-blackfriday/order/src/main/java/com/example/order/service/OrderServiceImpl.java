package com.example.order.service;

import com.example.order.domain.ProductOrder;
import com.example.order.model.OrderDTO;
import com.example.order.repos.ProductOrderRepository;
import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.NotNull;
import jakarta.validation.constraints.Positive;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Service;

import java.util.List;
import java.util.Map;
import java.util.UUID;

@Service
@Slf4j
@RequiredArgsConstructor
public class OrderServiceImpl implements OrderService {

    private final ProductOrderRepository productOrderRepository;

    @Override
    public OrderDTO.StartOrderResponse startOrder(OrderDTO.StartOrderRequest payload) {
        //TODO: check product in catalog-service
        //TODO: get payment
        //TODO: get delivery by user id
        ProductOrder order = new ProductOrder();
        order.setId(UUID.randomUUID().toString());
        order.setProductId("");
        order.setCount(payload.getCount());
        order.setPaymentId("");
        order.setDeliveryId("");
        order.setOrderStatus(ProductOrder.Status.CREATED);
        order = productOrderRepository.save(order);

        return OrderDTO.StartOrderResponse.builder()
                .orderId(order.getId())
                .paymentMethod(Map.of())
                .address(Map.of())
                .build();
    }

    @Override
    public OrderDTO.FinishOrderRequest finishOrder(OrderDTO.FinishOrderRequest payload) {
        return null;
    }

    @Override
    public OrderDTO.ProductDetail findById(String orderId) {
        return null;
    }

    @Override
    public List<OrderDTO.ProductOrder> findAllByUser(String userId) {
        return List.of();
    }
}
