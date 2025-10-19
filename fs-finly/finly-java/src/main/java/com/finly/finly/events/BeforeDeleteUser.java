package com.finly.finly.events;

import lombok.AllArgsConstructor;
import lombok.Getter;

import java.util.UUID;


@Getter
@AllArgsConstructor
public class BeforeDeleteUser {

    private UUID id;

}
