package com.example.timesale.rest.v2;

import com.example.timesale.domain.TimeSale;
import com.example.timesale.model.TimeSaleDTO;
import com.example.timesale.service.TimeSaleService;
import com.example.timesale.utils.UserIdInterceptor;
import jakarta.validation.Valid;
import org.springframework.beans.factory.annotation.Qualifier;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.data.web.PageableDefault;
import org.springframework.http.HttpStatus;
import org.springframework.http.MediaType;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

@RestController("TimeSaleResourceV2")
@RequestMapping(value = "/api/v2/timeSales", produces = MediaType.APPLICATION_JSON_VALUE)
public class TimeSaleResource {

    private final TimeSaleService timeSaleService;

    public TimeSaleResource(@Qualifier("TimeSaleServiceImplV2") TimeSaleService timeSaleService) {
        this.timeSaleService = timeSaleService;
    }

    @PostMapping
    public ResponseEntity<TimeSaleDTO.Response> createTimeSale(@RequestBody @Valid TimeSaleDTO.CreateRequest request) {
        TimeSale timeSale = timeSaleService.create(request);
        return ResponseEntity.status(HttpStatus.CREATED).body(TimeSaleDTO.Response.from(timeSale));
    }

    @GetMapping("/{timeSaleId}")
    public ResponseEntity<TimeSaleDTO.Response> getTimeSale(@PathVariable String timeSaleId) {
        TimeSale timeSale = timeSaleService.findById(timeSaleId);
        return ResponseEntity.ok(TimeSaleDTO.Response.from(timeSale));
    }

    @GetMapping
    public ResponseEntity<Page<TimeSaleDTO.Response>> getOngoingTimeSales(@PageableDefault Pageable pageable) {
        Page<TimeSale> timeSales = timeSaleService.findAllOngoing(pageable);
        return ResponseEntity.ok(timeSales.map(TimeSaleDTO.Response::from));
    }

    @PostMapping("/purchase")
    public ResponseEntity<TimeSaleDTO.PurchaseResponse> purchaseTimeSale(@RequestBody @Valid TimeSaleDTO.PurchaseRequest request) {
        String userId = UserIdInterceptor.getCurrentUserId();
        TimeSale timeSale = timeSaleService.purchase(request, userId);
        return ResponseEntity.status(HttpStatus.CREATED)
                .body(TimeSaleDTO.PurchaseResponse.from(timeSale, userId, request.getQuantity()));
    }
}
