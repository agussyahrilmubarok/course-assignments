package com.example.witrack.backend.service.impl;

import com.example.witrack.backend.common.BaseServiceTest;
import com.example.witrack.backend.domain.Ticket;
import com.example.witrack.backend.domain.User;
import com.example.witrack.backend.exception.UnauthorizedException;
import com.example.witrack.backend.model.TicketDTO;
import com.example.witrack.backend.repos.TicketRepository;
import com.example.witrack.backend.security.user.CurrentUserDetails;
import org.junit.jupiter.api.Assertions;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.Mockito;

import java.util.Collections;
import java.util.UUID;

class TicketServiceImplTest extends BaseServiceTest {

    @InjectMocks
    private TicketServiceImpl ticketService;

    @Mock
    private TicketRepository ticketRepository;

    @Mock
    private CurrentUserDetails currentUserDetails;

    private User testUser1;
    private Ticket testTicket1;
    private User testAdmin;
    private User testUser2;

    @BeforeEach
    void setUp() {
        testUser1 = new User();
        testUser1.setId(UUID.randomUUID());
        testUser1.setFullName("Jane Doe");
        testUser1.setEmail("janedoe@mail.com");
        testUser1.setPassword("encodedPassword");
        testUser1.setRoles(Collections.singleton(User.Role.ROLE_USER));

        testTicket1 = new Ticket();
        testTicket1.setId(UUID.randomUUID());
        testTicket1.setCode("TIC-1");
        testTicket1.setTitle("Failed Network Connection");
        testTicket1.setDescription("Low speed internet connection");
        testTicket1.setStatus(Ticket.Status.OPEN);
        testTicket1.setPriority(Ticket.Priority.LOW);
        testTicket1.setUser(testUser1);

        testAdmin = new User();
        testAdmin.setId(UUID.randomUUID());
        testAdmin.setFullName("Administrator");
        testAdmin.setEmail("admin@mail.com");
        testAdmin.setPassword("encodedPassword");
        testAdmin.setRoles(Collections.singleton(User.Role.ROLE_ADMIN));

        testUser2 = new User();
        testUser2.setId(UUID.randomUUID());
        testUser2.setFullName("John Doe");
        testUser2.setEmail("johndoe@mail.com");
        testUser2.setPassword("encodedPassword");
        testUser2.setRoles(Collections.singleton(User.Role.ROLE_USER));
    }

    @Test
    void givenValidRequest_whenCreateTicket_thenReturnTicketResponse() {
        Mockito.when(currentUserDetails.getUser()).thenReturn(testUser1);
        Mockito.when(ticketRepository.save(Mockito.any(Ticket.class))).thenReturn(testTicket1);

        TicketDTO.TicketRequest ticketRequest = TicketDTO.TicketRequest.builder()
                .title(testTicket1.getTitle())
                .description(testTicket1.getDescription())
                .status(testTicket1.getStatus().name())
                .priority(testTicket1.getPriority().name())
                .build();
        TicketDTO.TicketResponse response = ticketService.create(ticketRequest);

        Assertions.assertNotNull(response);
        Assertions.assertEquals(ticketRequest.getTitle(), response.getTitle());
        Assertions.assertEquals(ticketRequest.getDescription(), response.getDescription());
        Assertions.assertNotNull(response.getUser());
        Mockito.verify(ticketRepository, Mockito.times(1)).save(Mockito.any(Ticket.class));
    }

    @Test
    void givenSaveFails_whenCreateTicket_thenThrowException() {
        Mockito.when(currentUserDetails.getUser()).thenReturn(testUser1);
        Mockito.when(ticketRepository.save(Mockito.any(Ticket.class))).thenThrow(new RuntimeException("Database error"));

        TicketDTO.TicketRequest ticketRequest = TicketDTO.TicketRequest.builder()
                .title(testTicket1.getTitle())
                .description(testTicket1.getDescription())
                .status(testTicket1.getStatus().name())
                .priority(testTicket1.getPriority().name())
                .build();
        RuntimeException exception = Assertions.assertThrows(RuntimeException.class,
                () -> ticketService.create(ticketRequest));

        Assertions.assertEquals("Database error", exception.getMessage());
        Mockito.verify(ticketRepository, Mockito.times(1)).save(Mockito.any(Ticket.class));
    }

    @Test
    void givenOwner_whenUpdate_thenAccessGranted() {
        Mockito.when(currentUserDetails.getUser()).thenReturn(testUser1);
        Mockito.when(ticketRepository.findById(testTicket1.getId())).thenReturn(java.util.Optional.of(testTicket1));
        Mockito.when(ticketRepository.save(Mockito.any(Ticket.class))).thenReturn(testTicket1);

        TicketDTO.TicketRequest updateRequest = TicketDTO.TicketRequest.builder()
                .title("Updated Title")
                .build();
        TicketDTO.TicketResponse response = ticketService.update(testTicket1.getId(), updateRequest);

        Assertions.assertEquals("Updated Title", response.getTitle());
        Mockito.verify(ticketRepository, Mockito.times(1)).findById(Mockito.any());
        Mockito.verify(ticketRepository, Mockito.times(1)).save(Mockito.any(Ticket.class));
    }

    @Test
    void givenAdmin_whenUpdate_thenAccessGranted() {
        Mockito.when(currentUserDetails.getUser()).thenReturn(testAdmin);
        Mockito.when(ticketRepository.findById(testTicket1.getId())).thenReturn(java.util.Optional.of(testTicket1));
        Mockito.when(ticketRepository.save(Mockito.any(Ticket.class))).thenReturn(testTicket1);

        TicketDTO.TicketRequest updateRequest = TicketDTO.TicketRequest.builder()
                .title("Updated Title by Admin")
                .build();
        TicketDTO.TicketResponse response = ticketService.update(testTicket1.getId(), updateRequest);

        Assertions.assertEquals("Updated Title by Admin", response.getTitle());
        Mockito.verify(ticketRepository, Mockito.times(1)).findById(Mockito.any());
        Mockito.verify(ticketRepository, Mockito.times(1)).save(Mockito.any(Ticket.class));
    }

    @Test
    void givenNonOwnerNonAdmin_whenUpdate_thenThrowUnauthorized() {
        Mockito.when(currentUserDetails.getUser()).thenReturn(testUser2);
        Mockito.when(ticketRepository.findById(testTicket1.getId())).thenReturn(java.util.Optional.of(testTicket1));

        TicketDTO.TicketRequest updateRequest = TicketDTO.TicketRequest.builder()
                .title("Should Fail")
                .build();
        UnauthorizedException exception = Assertions.assertThrows(UnauthorizedException.class,
                () -> ticketService.update(testTicket1.getId(), updateRequest));

        Assertions.assertEquals("You do not have permission to access this ticket", exception.getMessage());
        Mockito.verify(ticketRepository, Mockito.never()).save(Mockito.any(Ticket.class));
    }

    @Test
    void givenOwner_whenDelete_thenSuccess() {
        Mockito.when(currentUserDetails.getUser()).thenReturn(testUser1);
        Mockito.when(ticketRepository.findById(testTicket1.getId())).thenReturn(java.util.Optional.of(testTicket1));

        Assertions.assertDoesNotThrow(() -> ticketService.delete(testTicket1.getId()));
        Mockito.verify(ticketRepository, Mockito.times(1)).delete(testTicket1);
    }

    @Test
    void givenNonOwnerNonAdmin_whenDelete_thenThrowUnauthorized() {
        Mockito.when(currentUserDetails.getUser()).thenReturn(testUser2);
        Mockito.when(ticketRepository.findById(testTicket1.getId())).thenReturn(java.util.Optional.of(testTicket1));

        UnauthorizedException exception = Assertions.assertThrows(UnauthorizedException.class,
                () -> ticketService.delete(testTicket1.getId()));

        Assertions.assertEquals("You do not have permission to access this ticket", exception.getMessage());
    }
}