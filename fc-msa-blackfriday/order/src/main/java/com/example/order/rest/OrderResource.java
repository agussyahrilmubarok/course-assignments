package com.example.order.rest;

import com.example.order.model.OrderDTO;
import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import org.springframework.http.MediaType;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.List;

@RestController("OrderResourceV1")
@RequestMapping(value = "/api/v1/orders", produces = MediaType.APPLICATION_JSON_VALUE)
@RequiredArgsConstructor
public class OrderResource {

    @PostMapping("/start")
    public ResponseEntity<OrderDTO.StartOrderResponse> startOrder(@RequestBody @Valid OrderDTO.StartOrderRequest payload) {
        throw new RuntimeException();
    }

    @PostMapping("/finish")
    public ResponseEntity<OrderDTO.FinishOrderRequest> finishOrder(@RequestBody @Valid OrderDTO.FinishOrderRequest payload) {
        throw new RuntimeException();
    }

    @GetMapping("/{orderId}")
    public ResponseEntity<OrderDTO.ProductDetail> findById(@PathVariable String orderId) {
        throw new RuntimeException();
    }

    @GetMapping("/users/{userId}")
    public ResponseEntity<List<OrderDTO.ProductOrder>> findAllByUser(@PathVariable String userId) {
        throw new RuntimeException();
    }
}
