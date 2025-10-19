package com.example.accounts.service.client;

import com.example.accounts.dto.CardsDto;
import org.springframework.stereotype.Component;

@Component
public class CardsFallback implements CardsFeignClient {

    @Override
    public CardsDto fetchCardDetails(String correlationId, String mobileNumber) {
        return null;
    }
}
