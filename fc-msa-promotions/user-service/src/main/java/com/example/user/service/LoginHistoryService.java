package com.example.user.service;

public interface LoginHistoryService {

    void recordLogin(String userId, String ipAddress);
}
