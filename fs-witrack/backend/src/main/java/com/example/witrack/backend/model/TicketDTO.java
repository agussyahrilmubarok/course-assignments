package com.example.witrack.backend.model;

import jakarta.validation.constraints.NotNull;
import jakarta.validation.constraints.Size;
import java.time.OffsetDateTime;
import java.util.UUID;
import lombok.Getter;
import lombok.Setter;


@Getter
@Setter
public class TicketDTO {

    private UUID id;

    @NotNull
    @Size(max = 255)
    @TicketCodeUnique
    private String code;

    @NotNull
    @Size(max = 255)
    private String title;

    private String description;

    @NotNull
    @Size(max = 255)
    private String status;

    @NotNull
    @Size(max = 255)
    private String priority;

    private OffsetDateTime completeAt;

    @NotNull
    private UUID user;

}
