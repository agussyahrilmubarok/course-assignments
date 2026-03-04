package com.example.witrack.backend.model;

import jakarta.validation.constraints.NotNull;
import java.util.UUID;
import lombok.Getter;
import lombok.Setter;


@Getter
@Setter
public class CommentDTO {

    private UUID id;

    private String content;

    @NotNull
    private UUID ticket;

    private UUID user;

}
