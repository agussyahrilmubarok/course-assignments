package com.example.point.rest.v1;

import com.example.point.model.PointDTO;
import com.example.point.service.v1.PointService;
import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.http.MediaType;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

@RestController("CouponPolicyResourceV1")
@RequestMapping(value = "/api/v1/points", produces = MediaType.APPLICATION_JSON_VALUE)
@RequiredArgsConstructor
public class PointResource {

    private final PointService pointService;

    @PostMapping("/earn")
    public ResponseEntity<PointDTO.Response> earnPoints(@RequestBody @Valid PointDTO.EarnRequest payload) {
        throw new RuntimeException();
    }

    @PostMapping("/use")
    public ResponseEntity<PointDTO.Response> usePoints(@RequestBody @Valid PointDTO.UseRequest payload) {
        throw new RuntimeException();
    }

    @PostMapping("/{pointId}/cancel")
    public ResponseEntity<PointDTO.Response> cancelPoints(@PathVariable Long pointId,
                                                          @RequestBody @Valid PointDTO.CancelRequest request) {
        throw new RuntimeException();
    }

    @GetMapping("/users/{userId}/balance")
    public ResponseEntity<PointDTO.BalanceResponse> getBalance(@PathVariable String userId) {
        throw new RuntimeException();
    }

    @GetMapping("/users/{userId}/history")
    public ResponseEntity<Page<PointDTO.Response>> getPointHistory(@PathVariable String userId,
                                                                   Pageable pageable) {
        throw new RuntimeException();
    }
}
