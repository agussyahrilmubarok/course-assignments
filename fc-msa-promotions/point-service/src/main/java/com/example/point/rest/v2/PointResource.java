package com.example.point.rest.v2;

import com.example.point.domain.Point;
import com.example.point.model.PointDTO;
import com.example.point.service.PointService;
import com.example.point.utils.UserIdInterceptor;
import io.swagger.v3.oas.annotations.Parameter;
import jakarta.validation.Valid;
import org.springframework.beans.factory.annotation.Qualifier;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.http.MediaType;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

@RestController("PointResourceV2")
@RequestMapping(value = "/api/v2/points", produces = MediaType.APPLICATION_JSON_VALUE)
public class PointResource {

    private final PointService pointService;

    public PointResource(@Qualifier("PointServiceImplV2") PointService pointService) {
        this.pointService = pointService;
    }

    @PostMapping("/earn")
    public ResponseEntity<PointDTO.Response> earnPoints(@RequestBody @Valid PointDTO.EarnRequest payload,
                                                        @Parameter(name = "X-USER-ID", description = "User ID", required = true)
                                                        @RequestHeader("X-USER-ID") String userId) {
        return ResponseEntity.ok(PointDTO.Response.from(pointService.earn(payload)));
    }

    @PostMapping("/use")
    public ResponseEntity<PointDTO.Response> usePoints(@RequestBody @Valid PointDTO.UseRequest payload,
                                                       @Parameter(name = "X-USER-ID", description = "User ID", required = true)
                                                       @RequestHeader("X-USER-ID") String userId) {
        return ResponseEntity.ok(PointDTO.Response.from(pointService.use(payload)));
    }

    @PostMapping("/cancel")
    public ResponseEntity<PointDTO.Response> cancelPoints(@RequestBody @Valid PointDTO.CancelRequest payload,
                                                          @Parameter(name = "X-USER-ID", description = "User ID", required = true)
                                                          @RequestHeader("X-USER-ID") String userId) {
        return ResponseEntity.ok(PointDTO.Response.from(pointService.cancel(payload)));
    }

    @GetMapping("/users/balance")
    public ResponseEntity<PointDTO.BalanceResponse> getBalance(@Parameter(name = "X-USER-ID", description = "User ID", required = true)
                                                               @RequestHeader("X-USER-ID") String userId) {
        return ResponseEntity.ok(PointDTO.BalanceResponse.of(UserIdInterceptor.getCurrentUserId(), pointService.getBalance()));
    }

    @GetMapping("/users/history")
    public ResponseEntity<Page<PointDTO.Response>> getPointHistory(Pageable pageable,
                                                                   @Parameter(name = "X-USER-ID", description = "User ID", required = true)
                                                                   @RequestHeader("X-USER-ID") String userId) {
        Page<Point> points = pointService.getHistory(pageable);
        return ResponseEntity.ok(points.map(PointDTO.Response::from));
    }
}
