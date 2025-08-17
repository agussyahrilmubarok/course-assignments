package com.example.witrack.backend.service.impl;

import com.example.witrack.backend.domain.Ticket;
import com.example.witrack.backend.domain.TicketReply;
import com.example.witrack.backend.domain.User;
import com.example.witrack.backend.exception.NotFoundException;
import com.example.witrack.backend.exception.UnauthorizedException;
import com.example.witrack.backend.model.*;
import com.example.witrack.backend.repository.TicketReplyRepository;
import com.example.witrack.backend.repository.TicketRepository;
import com.example.witrack.backend.repository.UserRepository;
import com.example.witrack.backend.security.CurrentUserDetails;
import com.example.witrack.backend.service.BaseServiceTest;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.mockito.InjectMocks;
import org.mockito.Mock;

import java.util.*;

import static org.junit.jupiter.api.Assertions.*;
import static org.mockito.Mockito.*;


class TicketServiceImplTest extends BaseServiceTest {

    @InjectMocks
    private TicketServiceImpl ticketService;

    @Mock
    private TicketRepository ticketRepository;

    @Mock
    private TicketReplyRepository ticketReplyRepository;

    @Mock
    private CurrentUserDetails currentUserDetails;

    @Mock
    private UserRepository userRepository;

    private User testUser1;
    private User testUser2;
    private User testAdmin;
    private Ticket testTicket1;
    private Ticket testTicket2;
    private TicketReply testTicketReply1;

    @BeforeEach
    void setUp() {
        testAdmin = new User();
        testAdmin.setId(UUID.randomUUID().toString());
        testAdmin.setFullName("John Doe");
        testAdmin.setEmail("johndoe@mail.com");
        testAdmin.setPassword("encodedPassword");
        testAdmin.setRoles(new HashSet<>(Arrays.asList(User.Role.ROLE_USER, User.Role.ROLE_ADMIN)));

        testUser1 = new User();
        testUser1.setId(UUID.randomUUID().toString());
        testUser1.setFullName("Jane Doe");
        testUser1.setEmail("janedoe@mail.com");
        testUser1.setPassword("encodedPassword");
        testUser1.setRoles(Collections.singleton(User.Role.ROLE_USER));

        testUser2 = new User();
        testUser2.setId(UUID.randomUUID().toString());
        testUser2.setFullName("Jane Doe");
        testUser2.setEmail("janedoe@mail.com");
        testUser2.setPassword("encodedPassword");
        testUser2.setRoles(Collections.singleton(User.Role.ROLE_USER));

        testTicket1 = new Ticket();
        testTicket1.setId(UUID.randomUUID().toString());
        testTicket1.setCode("TIC-1");
        testTicket1.setTitle("Ticket One");
        testTicket1.setDescription("Ticket One Description Network");
        testTicket1.setStatus(Ticket.Status.OPEN);
        testTicket1.setPriority(Ticket.Priority.LOW);
        testTicket1.setUser(testUser1);

        testTicketReply1 = new TicketReply();
        testTicketReply1.setId(UUID.randomUUID().toString());
        testTicketReply1.setContent("Content Ticket One");
        testTicketReply1.setTicket(testTicket1);
        testTicketReply1.setUser(testUser1);

        testTicket2 = new Ticket();
        testTicket2.setId(UUID.randomUUID().toString());
        testTicket2.setCode("TIC-2");
        testTicket2.setTitle("Ticket Two");
        testTicket2.setDescription("Ticket Two Description Cable");
        testTicket2.setStatus(Ticket.Status.RESOLVED);
        testTicket2.setPriority(Ticket.Priority.HIGH);
        testTicket2.setUser(testUser2);
    }

    @Test
    void givenNoFilter_whenGetTickets_thenReturnAllTickets() {
        when(ticketRepository.findAll()).thenReturn(Arrays.asList(testTicket1, testTicket2));

        List<TicketResponse> responses = ticketService.getTickets(null, null, null);

        assertEquals(2, responses.size());
        assertTrue(responses.stream().anyMatch(r -> r.getCode().equals(testTicket1.getCode())));

        verify(ticketRepository, times(1)).findAll();
    }

    @Test
    void givenSearchFilter_whenGetTickets_thenReturnFilteredTickets() {
        when(ticketRepository.findByCodeIgnoreCaseContainingOrTitleIgnoreCaseContainingOrDescriptionIgnoreCaseContaining(
                eq("network"), eq("network"), eq("network")
        )).thenReturn(Collections.singletonList(testTicket1));

        List<TicketResponse> responses = ticketService.getTickets("network", null, null);

        assertEquals(1, responses.size());
        assertEquals(testTicket1.getCode(), responses.get(0).getCode());

        verify(ticketRepository, times(1))
                .findByCodeIgnoreCaseContainingOrTitleIgnoreCaseContainingOrDescriptionIgnoreCaseContaining(
                        eq("network"), eq("network"), eq("network")
                );
    }

    @Test
    void givenStatusFilter_whenGetTickets_thenReturnFilteredTickets() {
        when(ticketRepository.findByStatus("OPEN")).thenReturn(Collections.singletonList(testTicket1));

        List<TicketResponse> responses = ticketService.getTickets(null, "OPEN", null);

        assertEquals(1, responses.size());
        assertEquals(testTicket1.getCode(), responses.get(0).getCode());

        verify(ticketRepository, times(1)).findByStatus("OPEN");
    }

    @Test
    void givenPriorityFilter_whenGetTickets_thenReturnFilteredTickets() {
        when(ticketRepository.findByPriority("LOW")).thenReturn(Collections.singletonList(testTicket1));

        List<TicketResponse> responses = ticketService.getTickets(null, null, "LOW");

        assertEquals(1, responses.size());
        assertEquals(testTicket1.getCode(), responses.get(0).getCode());

        verify(ticketRepository, times(1)).findByPriority("LOW");
    }

    @Test
    void givenStatusAndPriorityFilter_whenGetTickets_thenReturnFilteredTickets() {
        when(ticketRepository.findByStatusAndPriority("OPEN", "LOW"))
                .thenReturn(Collections.singletonList(testTicket1));

        List<TicketResponse> responses = ticketService.getTickets(null, "OPEN", "LOW");

        assertEquals(1, responses.size());
        assertEquals(testTicket1.getCode(), responses.get(0).getCode());

        verify(ticketRepository, times(1)).findByStatusAndPriority("OPEN", "LOW");
    }

    @Test
    void givenSearchAndStatusAndPriorityFilter_whenGetTickets_thenReturnEmpty() {
        when(ticketRepository.findByCodeIgnoreCaseContainingOrTitleIgnoreCaseContainingOrDescriptionIgnoreCaseContaining(
                eq("network"), eq("network"), eq("network")
        )).thenReturn(Collections.singletonList(testTicket1));

        when(ticketRepository.findByStatusAndPriority("OPEN", "HIGH")).thenReturn(Collections.emptyList());

        List<TicketResponse> responses = ticketService.getTickets("network", "OPEN", "HIGH");

        assertTrue(responses.isEmpty());

        verify(ticketRepository, times(1)).findByStatusAndPriority("OPEN", "HIGH");
    }

    @Test
    void givenExistingCode_whenGetTicketByCode_thenReturnTicketDetailResponse() {
        when(ticketRepository.findByCode(testTicket1.getCode())).thenReturn(Optional.of(testTicket1));

        TicketDetailResponse response = ticketService.getTicketByCode(testTicket1.getCode());

        assertNotNull(response);
        assertEquals(testTicket1.getCode(), response.getCode());
        assertEquals(testTicket1.getTitle(), response.getTitle());
        assertEquals(testTicket1.getDescription(), response.getDescription());

        verify(ticketRepository, times(1)).findByCode(testTicket1.getCode());
    }

    @Test
    void givenNonExistingCode_whenGetTicketByCode_thenThrowNotFoundException() {
        String nonExistingCode = "TIC-999";
        when(ticketRepository.findByCode(nonExistingCode)).thenReturn(Optional.empty());

        NotFoundException exception = assertThrows(
                NotFoundException.class,
                () -> ticketService.getTicketByCode(nonExistingCode)
        );

        assertTrue(exception.getMessage().contains(nonExistingCode));
        verify(ticketRepository, times(1)).findByCode(nonExistingCode);
    }

    @Test
    void givenValidRequest_whenCreateTicket_thenReturnTicketResponse() {
        TicketStoreRequest request = new TicketStoreRequest();
        request.setTitle("Test Ticket");
        request.setDescription("Test Description");
        request.setStatus("OPEN");
        request.setPriority("HIGH");

        when(currentUserDetails.getId()).thenReturn(testUser1.getId());
        when(userRepository.findById(testUser1.getId())).thenReturn(Optional.of(testUser1));
        when(ticketRepository.save(any(Ticket.class))).thenAnswer(invocation -> invocation.getArgument(0));

        TicketResponse response = ticketService.createTicket(request);

        assertNotNull(response);
        assertEquals(request.getTitle(), response.getTitle());
        assertEquals(request.getDescription(), response.getDescription());
        assertEquals(request.getStatus(), response.getStatus());
        assertEquals(request.getPriority(), response.getPriority());

        verify(ticketRepository, times(1)).save(any(Ticket.class));
    }

    @Test
    void givenOwnerUser_whenUpdateTicket_thenSuccess() {
        TicketStoreRequest request = new TicketStoreRequest();
        request.setTitle("Updated Title");

        testTicket1.setTitle(request.getTitle());
        when(ticketRepository.findByCode(testTicket1.getCode())).thenReturn(Optional.of(testTicket1));
        when(currentUserDetails.getId()).thenReturn(testUser1.getId());
        when(userRepository.findById(testUser1.getId())).thenReturn(Optional.of(testUser1));
        when(ticketRepository.save(any(Ticket.class))).thenAnswer(invocation -> invocation.getArgument(0));

        TicketResponse response = ticketService.updateTicket(testTicket1.getCode(), request);

        assertEquals("Updated Title", response.getTitle());
        verify(ticketRepository, times(1)).save(any(Ticket.class));
    }

    @Test
    void givenOwnerUser_whenUpdateMultipleFields_thenSuccess() {
        TicketStoreRequest request = new TicketStoreRequest();
        request.setTitle("New Title");
        request.setDescription("New Description");
        request.setStatus("ONPROGRESS");
        request.setPriority("CRITICAL");

        when(ticketRepository.findByCode(testTicket1.getCode())).thenReturn(Optional.of(testTicket1));
        when(currentUserDetails.getId()).thenReturn(testUser1.getId());
        when(userRepository.findById(testUser1.getId())).thenReturn(Optional.of(testUser1));
        when(ticketRepository.save(any(Ticket.class))).thenAnswer(invocation -> invocation.getArgument(0));

        TicketResponse response = ticketService.updateTicket(testTicket1.getCode(), request);

        assertEquals("New Title", response.getTitle());
        assertEquals("New Description", response.getDescription());
        assertEquals(Ticket.Status.ONPROGRESS.name(), response.getStatus());
        assertEquals(Ticket.Priority.CRITICAL.name(), response.getPriority());

        verify(ticketRepository, times(1)).save(any(Ticket.class));
    }

    @Test
    void givenNonOwnerButAdmin_whenUpdateTicket_thenSuccess() {
        TicketStoreRequest request = new TicketStoreRequest();
        request.setTitle("Admin Updated Title");

        User adminUser = new User();
        adminUser.setId(UUID.randomUUID().toString());
        adminUser.setRoles(Collections.singleton(User.Role.ROLE_ADMIN));

        when(ticketRepository.findByCode(testTicket1.getCode())).thenReturn(Optional.of(testTicket1));
        when(currentUserDetails.getId()).thenReturn(adminUser.getId());
        when(userRepository.findById(adminUser.getId())).thenReturn(Optional.of(adminUser));
        when(ticketRepository.save(any(Ticket.class))).thenAnswer(invocation -> invocation.getArgument(0));

        TicketResponse response = ticketService.updateTicket(testTicket1.getCode(), request);

        assertEquals("Admin Updated Title", response.getTitle());
        verify(ticketRepository, times(1)).save(any(Ticket.class));
    }

    @Test
    void givenNonOwnerNonAdmin_whenUpdateTicket_thenThrowUnauthorized() {
        TicketStoreRequest request = new TicketStoreRequest();
        request.setTitle("Hacked Title");

        User anotherUser = new User();
        anotherUser.setId(UUID.randomUUID().toString());
        anotherUser.setRoles(Collections.singleton(User.Role.ROLE_USER));

        when(ticketRepository.findByCode(testTicket1.getCode())).thenReturn(Optional.of(testTicket1));
        when(currentUserDetails.getId()).thenReturn(anotherUser.getId());
        when(userRepository.findById(anotherUser.getId())).thenReturn(Optional.of(anotherUser));

        assertThrows(UnauthorizedException.class,
                () -> ticketService.updateTicket(testTicket1.getCode(), request));

        verify(ticketRepository, never()).save(any(Ticket.class));
    }

    @Test
    void givenOwnerUser_whenDeleteTicket_thenSuccess() {
        when(ticketRepository.findByCode(testTicket1.getCode())).thenReturn(Optional.of(testTicket1));
        when(currentUserDetails.getId()).thenReturn(testUser1.getId());
        when(userRepository.findById(testUser1.getId())).thenReturn(Optional.of(testUser1));

        ticketService.deleteTicketByCode(testTicket1.getCode());

        verify(ticketRepository, times(1)).delete(testTicket1);
    }

    @Test
    void givenNonOwnerButAdmin_whenDeleteTicket_thenSuccess() {
        User adminUser = new User();
        adminUser.setId(UUID.randomUUID().toString());
        adminUser.setRoles(Collections.singleton(User.Role.ROLE_ADMIN));

        when(ticketRepository.findByCode(testTicket1.getCode())).thenReturn(Optional.of(testTicket1));
        when(currentUserDetails.getId()).thenReturn(adminUser.getId());
        when(userRepository.findById(adminUser.getId())).thenReturn(Optional.of(adminUser));

        ticketService.deleteTicketByCode(testTicket1.getCode());

        verify(ticketRepository, times(1)).delete(testTicket1);
    }

    @Test
    void givenNonOwnerNonAdmin_whenDeleteTicket_thenThrowUnauthorized() {
        User anotherUser = new User();
        anotherUser.setId(UUID.randomUUID().toString());
        anotherUser.setRoles(Collections.singleton(User.Role.ROLE_USER));

        when(ticketRepository.findByCode(testTicket1.getCode())).thenReturn(Optional.of(testTicket1));
        when(currentUserDetails.getId()).thenReturn(anotherUser.getId());
        when(userRepository.findById(anotherUser.getId())).thenReturn(Optional.of(anotherUser));

        assertThrows(UnauthorizedException.class,
                () -> ticketService.deleteTicketByCode(testTicket1.getCode()));

        verify(ticketRepository, never()).delete(any(Ticket.class));
    }

    @Test
    void givenValidRequest_whenCreateTicketReply_thenReturnReplyResponse() {
        TicketReplyStoreRequest request = new TicketReplyStoreRequest();
        request.setContent("This is a reply");

        when(ticketRepository.findByCode(testTicket1.getCode())).thenReturn(Optional.of(testTicket1));
        when(currentUserDetails.getId()).thenReturn(testUser1.getId());
        when(userRepository.findById(testUser1.getId())).thenReturn(Optional.of(testUser1));
        when(ticketReplyRepository.save(any(TicketReply.class))).thenAnswer(invocation -> invocation.getArgument(0));

        TicketReplyResponse response = ticketService.createTicketReply(testTicket1.getCode(), request);

        assertNotNull(response);
        assertEquals(request.getContent(), response.getContent());
        assertEquals(testUser1.getId(), response.getUser().getId());

        verify(ticketReplyRepository, times(1)).save(any(TicketReply.class));
    }
}