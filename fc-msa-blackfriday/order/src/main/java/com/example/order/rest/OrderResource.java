package com.example.order.rest;

import com.example.order.model.OrderDTO;
import com.example.order.service.OrderService;
import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import org.springframework.http.MediaType;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.List;

@RestController
@RequestMapping(value = "/api/v1/orders", produces = MediaType.APPLICATION_JSON_VALUE)
@RequiredArgsConstructor
public class OrderResource {

    private final OrderService orderService;

    @PostMapping("/start")
    public ResponseEntity<OrderDTO.StartOrderResponse> startOrder(@RequestBody @Valid OrderDTO.StartOrderRequest payload) {
        return ResponseEntity.ok(orderService.startOrder(payload));
    }

    @PostMapping("/finish")
    public ResponseEntity<OrderDTO.FinishOrderRequest> finishOrder(@RequestBody @Valid OrderDTO.FinishOrderRequest payload) {
        throw new RuntimeException();
    }

    @GetMapping("/{orderId}")
    public ResponseEntity<OrderDTO.Response> findById(@PathVariable String orderId) {
        throw new RuntimeException();
    }

    @GetMapping("/users/{userId}")
    public ResponseEntity<List<OrderDTO.Response>> findAllByUser(@PathVariable String userId) {
        throw new RuntimeException();
    }
}
