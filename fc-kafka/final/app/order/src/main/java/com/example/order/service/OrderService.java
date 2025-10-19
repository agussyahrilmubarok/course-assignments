package com.example.order.service;

import com.example.order.domain.Order;
import com.example.order.kafka.OrderEvent;
import com.example.order.kafka.OrderEventProducer;
import com.example.order.model.OrderDTO;
import com.example.order.repos.OrderRepository;
import com.example.order.util.NotFoundException;
import com.example.order.util.OrderEventPublishException;
import com.fasterxml.jackson.core.JsonProcessingException;
import org.springframework.data.domain.Sort;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.time.OffsetDateTime;
import java.time.ZoneOffset;
import java.util.List;
import java.util.UUID;


@Service
public class OrderService {

    private static final ZoneOffset DEFAULT_ZONE_OFFSET = ZoneOffset.ofHours(7); // Asia/Jakarta

    private final OrderRepository orderRepository;
    private final OrderEventProducer orderEventProducer;

    public OrderService(final OrderRepository orderRepository, OrderEventProducer orderEventProducer) {
        this.orderRepository = orderRepository;
        this.orderEventProducer = orderEventProducer;
    }

    public List<OrderDTO> findAll() {
        final List<Order> orders = orderRepository.findAll(Sort.by("id"));
        return orders.stream()
                .map(order -> mapToDTO(order, new OrderDTO()))
                .toList();
    }

    public OrderDTO get(final UUID id) {
        return orderRepository.findById(id)
                .map(order -> mapToDTO(order, new OrderDTO()))
                .orElseThrow(NotFoundException::new);
    }

    @Transactional
    public UUID create(final OrderDTO orderDTO) {
        final Order order = new Order();
        mapToEntity(orderDTO, order);
        if (order.getStatus() == null) {
            order.setStatus(Order.Status.CREATED);
        }
        final Order saved = orderRepository.save(order);

        final OrderEvent orderEvent = OrderEvent.builder()
                .eventId(UUID.randomUUID())
                .eventType(OrderEvent.OrderEventType.CREATED)
                .eventAt(OffsetDateTime.now())
                .order(OrderEvent.OrderPayload.from(saved))
                .build();

        try {
            orderEventProducer.sendEvent(orderEvent);
        } catch (JsonProcessingException exception) {
            throw new OrderEventPublishException("Failed to publish OrderEvent to Kafka", exception);
        }

        return saved.getId();
    }

    public void update(final UUID id, final OrderDTO orderDTO) {
        final Order order = orderRepository.findById(id)
                .orElseThrow(NotFoundException::new);
        mapToEntity(orderDTO, order);
        orderRepository.save(order);
    }

    public void delete(final UUID id) {
        final Order order = orderRepository.findById(id)
                .orElseThrow(NotFoundException::new);
        orderRepository.delete(order);
    }

    private OrderDTO mapToDTO(final Order order, final OrderDTO orderDTO) {
        orderDTO.setId(order.getId());
        orderDTO.setCustomerId(order.getCustomerId());
        orderDTO.setProductId(order.getProductId());
        orderDTO.setQuantity(order.getQuantity());
        orderDTO.setTotalAmount(order.getTotalAmount());
        orderDTO.setStatus(order.getStatus() != null ? order.getStatus().name() : null);
        orderDTO.setOrderAt(order.getOrderAt() != null ? order.getOrderAt().toLocalDateTime() : null);
        return orderDTO;
    }

    private Order mapToEntity(final OrderDTO orderDTO, final Order order) {
        order.setCustomerId(orderDTO.getCustomerId());
        order.setProductId(orderDTO.getProductId());
        order.setQuantity(orderDTO.getQuantity());
        order.setTotalAmount(orderDTO.getTotalAmount());
        if (orderDTO.getStatus() != null) {
            order.setStatus(Order.Status.valueOf(orderDTO.getStatus()));
        }
        order.setOrderAt(orderDTO.getOrderAt() != null ? orderDTO.getOrderAt().atOffset(DEFAULT_ZONE_OFFSET) : null);
        return order;
    }
}
