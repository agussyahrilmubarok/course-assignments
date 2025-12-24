package com.example.witrack.backend.service.impl;

import com.example.witrack.backend.common.BaseServiceTest;
import com.example.witrack.backend.domain.Ticket;
import com.example.witrack.backend.domain.TicketComment;
import com.example.witrack.backend.domain.User;
import com.example.witrack.backend.exception.NotFoundException;
import com.example.witrack.backend.model.TicketCommentDTO;
import com.example.witrack.backend.repos.TicketCommentRepository;
import com.example.witrack.backend.repos.TicketRepository;
import com.example.witrack.backend.security.user.CurrentUserDetails;
import org.junit.jupiter.api.Assertions;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.Mockito;

import java.util.Collections;
import java.util.Optional;
import java.util.UUID;

class TicketCommentServiceImplTest extends BaseServiceTest {

    @InjectMocks
    private TicketCommentServiceImpl ticketCommentService;

    @Mock
    private TicketCommentRepository ticketCommentRepository;

    @Mock
    private TicketRepository ticketRepository;

    @Mock
    private CurrentUserDetails currentUserDetails;

    private User testUser;
    private Ticket testTicket;
    private TicketComment testComment;

    @BeforeEach
    void setUp() {
        testUser = new User();
        testUser.setId(UUID.randomUUID());
        testUser.setFullName("Jane Doe");
        testUser.setEmail("janedoe@mail.com");
        testUser.setPassword("encodedPassword");
        testUser.setRoles(Collections.singleton(User.Role.ROLE_USER));

        testTicket = new Ticket();
        testTicket.setId(UUID.randomUUID());
        testTicket.setCode("TIC-1");
        testTicket.setTitle("Failed Network Connection");
        testTicket.setDescription("Low speed internet connection");
        testTicket.setStatus(Ticket.Status.OPEN);
        testTicket.setPriority(Ticket.Priority.LOW);
        testTicket.setUser(testUser);

        testComment = new TicketComment();
        testComment.setId(UUID.randomUUID());
        testComment.setContent("This is a comment");
        testComment.setUser(testUser);
        testComment.setTicket(testTicket);
    }

    @Test
    void givenValidTicketCode_whenCreateComment_thenReturnCommentResponse() {
        Mockito.when(ticketRepository.findByCode(testTicket.getCode())).thenReturn(Optional.of(testTicket));
        Mockito.when(currentUserDetails.getUser()).thenReturn(testUser);
        Mockito.when(ticketCommentRepository.save(Mockito.any(TicketComment.class))).thenReturn(testComment);

        TicketCommentDTO.TicketCommentRequest request = TicketCommentDTO.TicketCommentRequest.builder()
                .content("This is a comment")
                .build();
        TicketCommentDTO.TicketCommentResponse response = ticketCommentService.create(testTicket.getCode(), request);

        Assertions.assertNotNull(response);
        Assertions.assertEquals(testComment.getContent(), response.getContent());
        Assertions.assertEquals(testUser.getId().toString(), response.getUser().getId());
        Mockito.verify(ticketRepository, Mockito.times(1)).findByCode(testTicket.getCode());
        Mockito.verify(ticketCommentRepository, Mockito.times(1)).save(Mockito.any(TicketComment.class));
    }

    @Test
    void givenInvalidTicketCode_whenCreateComment_thenThrowNotFoundException() {
        String invalidCode = "TIC-NOT-EXIST";
        Mockito.when(ticketRepository.findByCode(invalidCode)).thenReturn(Optional.empty());

        TicketCommentDTO.TicketCommentRequest request = TicketCommentDTO.TicketCommentRequest.builder()
                .content("This is a comment")
                .build();
        NotFoundException exception = Assertions.assertThrows(NotFoundException.class,
                () -> ticketCommentService.create(invalidCode, request));

        Assertions.assertEquals("Ticket not found", exception.getMessage());
        Mockito.verify(ticketRepository, Mockito.times(1)).findByCode(invalidCode);
        Mockito.verify(ticketCommentRepository, Mockito.never()).save(Mockito.any(TicketComment.class));
    }
}