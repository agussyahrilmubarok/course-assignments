package com.example.user.model;


import lombok.Data;

import java.io.Serializable;
import java.time.LocalDateTime;

@Data
public class LoginHistory implements Serializable {

    private String userId;

    private String ipAddress;

    private LocalDateTime loginAt;
}
