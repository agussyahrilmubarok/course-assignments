package com.finly.finly.model;

import lombok.Getter;
import lombok.Setter;

@Getter
@Setter
public class InfoMessage {

    private String message;
    private String type = "success";
}
