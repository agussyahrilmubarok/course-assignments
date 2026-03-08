package com.example.witrack.backend.rest.v1;

import com.example.witrack.backend.model.*;
import com.example.witrack.backend.rest.BaseResourceTest;
import com.example.witrack.backend.service.impl.TicketServiceImpl;
import lombok.SneakyThrows;
import org.junit.jupiter.api.Test;
import org.mockito.Mockito;
import org.springframework.http.MediaType;
import org.springframework.security.test.context.support.WithMockUser;
import org.springframework.test.context.bean.override.mockito.MockitoBean;
import org.springframework.test.web.servlet.request.MockMvcRequestBuilders;
import org.springframework.test.web.servlet.result.MockMvcResultMatchers;

import java.time.OffsetDateTime;
import java.util.List;
import java.util.UUID;

import static org.mockito.Mockito.*;

class TicketResourceTest extends BaseResourceTest {

    @MockitoBean
    private TicketServiceImpl ticketService;

    @Test
    @WithMockUser(username = "johndoe@mail.com", roles = {"USER"})
    @SneakyThrows
    void givenValidRequest_whenGetTickets_thenReturnList() {
        UserResponse userResponse = UserResponse.fromUser(mockUser);
        TicketResponse response = TicketResponse.builder()
                .id(UUID.randomUUID().toString())
                .code("TCK-001")
                .title("Sample Ticket")
                .status("OPEN")
                .priority("HIGH")
                .createdAt(OffsetDateTime.now())
                .user(userResponse)
                .build();

        Mockito.when(ticketService.getTickets(null, null, null, null)).thenReturn(List.of(response));

        mockMvc.perform(MockMvcRequestBuilders.get("/api/v1/tickets")
                        .contentType(MediaType.APPLICATION_JSON))
                .andExpect(MockMvcResultMatchers.status().isOk())
                .andExpect(MockMvcResultMatchers.jsonPath("$[0].title").value(response.getTitle()));

        verify(ticketService, times(1)).getTickets(null, null, null, null);
    }

    @Test
    @WithMockUser(username = "janedoe@mail.com", roles = {"USER"})
    @SneakyThrows
    void givenValidRequest_whenGetMyTickets_thenReturnUserTickets() {
        UserResponse userResponse = UserResponse.fromUser(mockUser);
        TicketResponse response = TicketResponse.builder()
                .id(UUID.randomUUID().toString())
                .code("TCK-002")
                .title("My Ticket")
                .status("RESOLVED")
                .priority("LOW")
                .createdAt(OffsetDateTime.now())
                .user(userResponse)
                .build();

        when(ticketService.getMyTickets(null, null, null, null)).thenReturn(List.of(response));

        mockMvc.perform(MockMvcRequestBuilders.get("/api/v1/tickets/me")
                        .contentType(MediaType.APPLICATION_JSON))
                .andExpect(MockMvcResultMatchers.status().isOk())
                .andExpect(MockMvcResultMatchers.jsonPath("$[0].title").value(response.getTitle()))
                .andExpect(MockMvcResultMatchers.jsonPath("$[0].code").value(response.getCode()));

        verify(ticketService, times(1)).getMyTickets(null, null, null, null);
    }

    @Test
    @WithMockUser(username = "johndoe@mail.com", roles = {"USER"})
    @SneakyThrows
    void givenValidCode_whenGetTicketByCode_thenReturnDetail() {
        UserResponse userResponse = UserResponse.fromUser(mockUser);
        TicketReplyResponse ticketReplyResponse = TicketReplyResponse.builder()
                .id(UUID.randomUUID().toString())
                .content("Content Ticket Reply")
                .user(userResponse)
                .build();
        TicketDetailResponse response = TicketDetailResponse.builder()
                .id(UUID.randomUUID().toString())
                .code("TCK-001")
                .title("Sample Ticket")
                .description("Sample Ticket Description")
                .status("OPEN")
                .priority("HIGH")
                .createdAt(OffsetDateTime.now())
                .user(userResponse)
                .replies(List.of(ticketReplyResponse))
                .build();

        Mockito.when(ticketService.getTicketByCode("TCK-001")).thenReturn(response);

        mockMvc.perform(MockMvcRequestBuilders.get("/api/v1/tickets/TCK-001")
                        .contentType(MediaType.APPLICATION_JSON))
                .andExpect(MockMvcResultMatchers.status().isOk())
                .andExpect(MockMvcResultMatchers.jsonPath("$.title").value(response.getTitle()));

        verify(ticketService, times(1)).getTicketByCode("TCK-001");
    }

    @Test
    @WithMockUser(username = "johndoe@mail.com", roles = {"USER"})
    @SneakyThrows
    void givenValidRequest_whenCreateTicket_thenReturnCreatedTicket() {
        TicketStoreRequest request = new TicketStoreRequest();
        request.setTitle("New Ticket");
        request.setDescription("New Desc");
        request.setStatus("OPEN");
        request.setPriority("MEDIUM");

        UserResponse userResponse = UserResponse.fromUser(mockUser);
        TicketResponse response = TicketResponse.builder()
                .id(UUID.randomUUID().toString())
                .code("TCK-003")
                .title("New Ticket")
                .description("New Desc")
                .status("OPEN")
                .priority("MEDIUM")
                .createdAt(OffsetDateTime.now())
                .user(userResponse)
                .build();

        when(ticketService.createTicket(any(TicketStoreRequest.class))).thenReturn(response);

        mockMvc.perform(MockMvcRequestBuilders.post("/api/v1/tickets")
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(objectMapper.writeValueAsString(request)))
                .andExpect(MockMvcResultMatchers.status().isOk())
                .andExpect(MockMvcResultMatchers.jsonPath("$.title").value("New Ticket"));

        verify(ticketService, times(1)).createTicket(any(TicketStoreRequest.class));
    }

    @Test
    @WithMockUser(username = "johndoe@mail.com", roles = {"USER", "ADMIN"})
    @SneakyThrows
    void givenValidRequest_whenUpdateTicket_thenReturnUpdatedTicket() {
        TicketStoreRequest request = new TicketStoreRequest();
        request.setTitle("Updated Ticket");
        request.setDescription("New Desc");
        request.setStatus("OPEN");
        request.setPriority("HIGH");

        UserResponse userResponse = UserResponse.fromUser(mockUser);
        TicketResponse response = TicketResponse.builder()
                .id(UUID.randomUUID().toString())
                .code("TCK-004")
                .title("Updated Ticket")
                .description("New Desc")
                .status("OPEN")
                .priority("HIGH")
                .createdAt(OffsetDateTime.now())
                .updatedAt(OffsetDateTime.now())
                .user(userResponse)
                .build();

        when(ticketService.updateTicket(eq("TCK-004"), any(TicketStoreRequest.class))).thenReturn(response);

        mockMvc.perform(MockMvcRequestBuilders.put("/api/v1/tickets/TCK-004")
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(objectMapper.writeValueAsString(request)))
                .andExpect(MockMvcResultMatchers.status().isOk())
                .andExpect(MockMvcResultMatchers.jsonPath("$.title").value("Updated Ticket"));

        verify(ticketService, times(1)).updateTicket(eq("TCK-004"), any(TicketStoreRequest.class));
    }

    @Test
    @WithMockUser(username = "johndoe@mail.com", roles = {"USER", "ADMIN"})
    @SneakyThrows
    void givenValidCode_whenDeleteTicket_thenReturnNoContent() {
        doNothing().when(ticketService).deleteTicketByCode("TCK-005");

        mockMvc.perform(MockMvcRequestBuilders.delete("/api/v1/tickets/TCK-005"))
                .andExpect(MockMvcResultMatchers.status().isNoContent());

        verify(ticketService, times(1)).deleteTicketByCode("TCK-005");
    }

    @Test
    @WithMockUser(username = "johndoe@mail.com", roles = {"USER"})
    @SneakyThrows
    void givenValidRequest_whenCreateReply_thenReturnReply() throws Exception {
        TicketReplyStoreRequest request = new TicketReplyStoreRequest();
        request.setContent("Reply Message");

        TicketReplyResponse response = TicketReplyResponse.builder()
                .id(UUID.randomUUID().toString())
                .content("Reply Message")
                .createdAt(OffsetDateTime.now())
                .build();

        when(ticketService.createTicketReply(eq("TCK-006"), any(TicketReplyStoreRequest.class))).thenReturn(response);

        mockMvc.perform(MockMvcRequestBuilders.post("/api/v1/tickets/TCK-006/reply")
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(objectMapper.writeValueAsString(request)))
                .andExpect(MockMvcResultMatchers.status().isOk());

        verify(ticketService, times(1)).createTicketReply(eq("TCK-006"), any(TicketReplyStoreRequest.class));
    }
}