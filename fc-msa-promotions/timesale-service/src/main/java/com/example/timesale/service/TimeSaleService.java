package com.example.timesale.service;

import com.example.timesale.domain.TimeSale;
import com.example.timesale.model.TimeSaleDTO;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;

public interface TimeSaleService {

    TimeSale create(TimeSaleDTO.CreateRequest request);

    TimeSale findById(String timeSaleId);

    Page<TimeSale> findAllOngoing(Pageable pageable);

    TimeSale purchase(TimeSaleDTO.PurchaseRequest request, String userId);
}
