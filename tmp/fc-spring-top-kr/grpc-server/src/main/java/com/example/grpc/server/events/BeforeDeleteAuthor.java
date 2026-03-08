package com.example.grpc.server.events;

import lombok.AllArgsConstructor;
import lombok.Getter;


@Getter
@AllArgsConstructor
public class BeforeDeleteAuthor {

    private Long id;

}
