package com.example.coupon.rest;

import com.example.coupon.domain.Coupon;
import com.example.coupon.model.CouponDTO;
import com.example.coupon.service.CouponService;
import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import org.springframework.http.HttpStatus;
import org.springframework.http.MediaType;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.List;

@RestController
@RequestMapping(value = "/api/coupons", produces = MediaType.APPLICATION_JSON_VALUE)
@RequiredArgsConstructor
public class CouponResource {

    private final CouponService couponService;

    @GetMapping
    public ResponseEntity<List<CouponDTO.Response>> getCoupons(@RequestParam(required = false) Coupon.Status status,
                                                               @RequestParam(required = false) Integer page,
                                                               @RequestParam(required = false) Integer size) {
        CouponDTO.ListRequest request = CouponDTO.ListRequest.builder()
                .status(status).page(page).size(size).build();
        List<CouponDTO.Response> responses = couponService.findCoupons(request);
        return ResponseEntity.ok(responses);
    }

    @GetMapping("/{couponId}")
    public ResponseEntity<CouponDTO.Response> getCoupon(@PathVariable String couponId) {
        CouponDTO.Response response = couponService.findCoupon(couponId);
        return new ResponseEntity<>(response, HttpStatus.OK);
    }

    @PostMapping("/issue")
    public ResponseEntity<CouponDTO.Response> issueCoupon(@RequestBody @Valid CouponDTO.IssueRequest payload) {
        CouponDTO.Response response = CouponDTO.Response.from(couponService.issueCoupon(payload));
        return new ResponseEntity<>(response, HttpStatus.CREATED);
    }

    @PostMapping("/{couponId}/use")
    public ResponseEntity<CouponDTO.Response> useCoupon(@PathVariable String couponId,
                                                        @RequestBody CouponDTO.UseRequest payload) {
        CouponDTO.Response response = CouponDTO.Response.from(couponService.useCoupon(couponId, payload.getOrderId()));
        return new ResponseEntity<>(response, HttpStatus.OK);
    }

    @PostMapping("/{couponId}/cancel")
    public ResponseEntity<CouponDTO.Response> cancelCoupon(@PathVariable String couponId) {
        CouponDTO.Response response = CouponDTO.Response.from(couponService.cancelCoupon(couponId));
        return new ResponseEntity<>(response, HttpStatus.OK);
    }
}
