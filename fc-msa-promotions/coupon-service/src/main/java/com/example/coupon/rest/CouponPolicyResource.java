package com.example.coupon.rest;

import com.example.coupon.model.CouponPolicyDTO;
import com.example.coupon.service.CouponPolicyService;
import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import org.springframework.http.HttpStatus;
import org.springframework.http.MediaType;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.List;
import java.util.stream.Collectors;

@RestController
@RequestMapping(value = "/api/couponPolicies", produces = MediaType.APPLICATION_JSON_VALUE)
@RequiredArgsConstructor
public class CouponPolicyResource {

    private final CouponPolicyService couponPolicyService;

    @GetMapping
    public ResponseEntity<List<CouponPolicyDTO.Response>> getCouponPolicies() {
        List<CouponPolicyDTO.Response> response = couponPolicyService.findAll().stream()
                .map(CouponPolicyDTO.Response::from).collect(Collectors.toList());
        return ResponseEntity.ok(response);
    }

    @GetMapping("/{id}")
    public ResponseEntity<CouponPolicyDTO.Response> getCouponPolicy(@PathVariable String id) {
        CouponPolicyDTO.Response response = CouponPolicyDTO.Response.from(couponPolicyService.findById(id));
        return new ResponseEntity<>(response, HttpStatus.OK);
    }

    @PostMapping
    public ResponseEntity<CouponPolicyDTO.Response> createCouponPolicy(@RequestBody @Valid CouponPolicyDTO.CreateRequest payload) {
        CouponPolicyDTO.Response response = CouponPolicyDTO.Response.from(couponPolicyService.create(payload));
        return new ResponseEntity<>(response, HttpStatus.CREATED);
    }
}
