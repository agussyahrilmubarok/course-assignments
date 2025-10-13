package com.example.timesale.rest.v3;

import com.example.timesale.domain.TimeSale;
import com.example.timesale.model.TimeSaleDTO;
import com.example.timesale.service.TimeSaleAsyncService;
import com.example.timesale.service.TimeSaleService;
import com.example.timesale.utils.UserIdInterceptor;
import io.swagger.v3.oas.annotations.Parameter;
import jakarta.validation.Valid;
import org.springframework.beans.factory.annotation.Qualifier;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.data.web.PageableDefault;
import org.springframework.http.HttpStatus;
import org.springframework.http.MediaType;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

@RestController("TimeSaleResourceV3")
@RequestMapping(value = "/api/v3/timeSales", produces = MediaType.APPLICATION_JSON_VALUE)
public class TimeSaleResource {

    private final TimeSaleService timeSaleService;
    private final TimeSaleAsyncService timeSaleAsyncService;

    public TimeSaleResource(@Qualifier("TimeSaleServiceImplV3") TimeSaleService timeSaleService,
                            @Qualifier("TimeSaleServiceImplV3") TimeSaleAsyncService timeSaleAsyncService) {
        this.timeSaleService = timeSaleService;
        this.timeSaleAsyncService = timeSaleAsyncService;
    }

    @PostMapping
    public ResponseEntity<TimeSaleDTO.Response> createTimeSale(@RequestBody @Valid TimeSaleDTO.CreateRequest request,
                                                               @Parameter(name = "X-USER-ID", description = "User ID", required = true)
                                                               @RequestHeader("X-USER-ID") String userIdSwagger) {
        TimeSale timeSale = timeSaleService.create(request);
        return ResponseEntity.status(HttpStatus.CREATED).body(TimeSaleDTO.Response.from(timeSale));
    }

    @GetMapping("/{timeSaleId}")
    public ResponseEntity<TimeSaleDTO.Response> getTimeSale(@PathVariable String timeSaleId,
                                                            @Parameter(name = "X-USER-ID", description = "User ID", required = true)
                                                            @RequestHeader("X-USER-ID") String userIdSwagger) {
        TimeSale timeSale = timeSaleService.findById(timeSaleId);
        return ResponseEntity.ok(TimeSaleDTO.Response.from(timeSale));
    }

    @GetMapping
    public ResponseEntity<Page<TimeSaleDTO.Response>> getOngoingTimeSales(@PageableDefault Pageable pageable,
                                                                          @Parameter(name = "X-USER-ID", description = "User ID", required = true)
                                                                          @RequestHeader("X-USER-ID") String userIdSwagger) {
        Page<TimeSale> timeSales = timeSaleService.findAllOngoing(pageable);
        return ResponseEntity.ok(timeSales.map(TimeSaleDTO.Response::from));
    }

    @PostMapping("/purchase")
    public ResponseEntity<TimeSaleDTO.AsyncPurchaseResponse> purchaseTimeSale(@RequestBody @Valid TimeSaleDTO.PurchaseRequest request,
                                                                              @Parameter(name = "X-USER-ID", description = "User ID", required = true)
                                                                              @RequestHeader("X-USER-ID") String userIdSwagger) {
        String userId = UserIdInterceptor.getCurrentUserId();
        String requestId = timeSaleAsyncService.purchaseRequest(request, userId);
        return ResponseEntity.ok(TimeSaleDTO.AsyncPurchaseResponse.builder()
                .requestId(requestId)
                .status("PENDING")
                .build());
    }

    @GetMapping("/{timeSaleId}/purchase/results/{requestId}")
    public ResponseEntity<TimeSaleDTO.AsyncPurchaseResponse> getPurchaseResult(@PathVariable String timeSaleId,
                                                                               @PathVariable String requestId,
                                                                               @Parameter(name = "X-USER-ID", description = "User ID", required = true)
                                                                               @RequestHeader("X-USER-ID") String userIdSwagger) {
        return ResponseEntity.ok(timeSaleAsyncService.findPurchaseResult(timeSaleId, requestId));
    }
}
