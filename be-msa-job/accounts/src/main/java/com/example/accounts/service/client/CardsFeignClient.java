package com.example.accounts.service.client;

import com.example.accounts.dto.CardsDto;
import org.springframework.cloud.openfeign.FeignClient;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestHeader;
import org.springframework.web.bind.annotation.RequestParam;

@FeignClient(name = "cards", path = "/api", fallback = CardsFallback.class)
public interface CardsFeignClient {

    @GetMapping("/fetch")
    CardsDto fetchCardDetails(@RequestHeader("examplebank-correlation-id") String correlationId,
                              @RequestParam String mobileNumber);

}
