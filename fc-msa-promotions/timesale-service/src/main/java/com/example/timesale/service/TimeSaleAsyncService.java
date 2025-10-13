package com.example.timesale.service;

import com.example.timesale.model.TimeSaleDTO;

public interface TimeSaleAsyncService {

    String purchaseRequest(TimeSaleDTO.PurchaseRequest request, String userId);

    TimeSaleDTO.AsyncPurchaseResponse findPurchaseResult(String timeSaleId, String requestId);

    void savePurchaseResult(String requestId, String result);

    void removePurchaseResultFromQueue(String timeSaleId, String requestId);
}
