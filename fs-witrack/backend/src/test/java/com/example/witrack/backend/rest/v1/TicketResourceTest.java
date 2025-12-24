package com.example.witrack.backend.rest.v1;

import com.example.witrack.backend.common.BaseResourceTest;
import com.example.witrack.backend.model.TicketCommentDTO;
import com.example.witrack.backend.model.TicketDTO;
import com.example.witrack.backend.model.UserDTO;
import com.example.witrack.backend.service.v1.impl.TicketCommentServiceImpl;
import com.example.witrack.backend.service.v1.impl.TicketServiceImpl;
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

class TicketResourceTest extends BaseResourceTest {

    @MockitoBean
    private TicketServiceImpl ticketService;

    @MockitoBean
    private TicketCommentServiceImpl ticketCommentService;

    @Test
    @WithMockUser(username = "johndoe@mail.com", roles = {"USER"})
    @SneakyThrows
    void givenValidRequest_whenCreateTicket_thenReturnCreatedTicket() {
        TicketDTO.TicketRequest request = TicketDTO.TicketRequest.builder()
                .title("New Ticket")
                .description("New Desc")
                .status("OPEN")
                .priority("MEDIUM")
                .build();
        UserDTO.UserResponse userResponse = UserDTO.UserResponse.fromUser(mockUser);
        TicketDTO.TicketResponse response = TicketDTO.TicketResponse.builder()
                .id(UUID.randomUUID().toString())
                .code("TCK-003")
                .title("New Ticket")
                .description("New Desc")
                .status("OPEN")
                .priority("MEDIUM")
                .createdAt(OffsetDateTime.now())
                .user(userResponse)
                .build();
        Mockito.when(ticketService.create(Mockito.any(TicketDTO.TicketRequest.class))).thenReturn(response);

        mockMvc.perform(MockMvcRequestBuilders.post("/api/v1/tickets")
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(objectMapper.writeValueAsString(request)))
                .andExpect(MockMvcResultMatchers.status().isOk())
                .andExpect(MockMvcResultMatchers.jsonPath("$.title").value("New Ticket"));

        Mockito.verify(ticketService, Mockito.times(1)).create(Mockito.any(TicketDTO.TicketRequest.class));
    }

    @Test
    @WithMockUser(username = "johndoe@mail.com", roles = {"USER", "ADMIN"})
    @SneakyThrows
    void givenValidRequest_whenUpdateTicket_thenReturnUpdatedTicket() {
        UUID id = UUID.randomUUID();
        TicketDTO.TicketRequest request = TicketDTO.TicketRequest.builder()
                .title("New Ticket Update")
                .description("New Desc")
                .status("OPEN")
                .priority("MEDIUM")
                .build();
        UserDTO.UserResponse userResponse = UserDTO.UserResponse.fromUser(mockUser);
        TicketDTO.TicketResponse response = TicketDTO.TicketResponse.builder()
                .id(id.toString())
                .code("TCK-003")
                .title("New Ticket Update")
                .description("New Desc")
                .status("OPEN")
                .priority("MEDIUM")
                .createdAt(OffsetDateTime.now())
                .user(userResponse)
                .build();
        Mockito.when(ticketService.update(Mockito.eq(id), Mockito.any(TicketDTO.TicketRequest.class))).thenReturn(response);

        mockMvc.perform(MockMvcRequestBuilders.put("/api/v1/tickets/{id}", id)
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(objectMapper.writeValueAsString(request)))
                .andExpect(MockMvcResultMatchers.status().isOk())
                .andExpect(MockMvcResultMatchers.jsonPath("$.title").value("New Ticket Update"));

        Mockito.verify(ticketService, Mockito.times(1)).update(Mockito.eq(id), Mockito.any(TicketDTO.TicketRequest.class));
    }

    @Test
    @WithMockUser(username = "johndoe@mail.com", roles = {"USER", "ADMIN"})
    @SneakyThrows
    void givenValidCode_whenDeleteTicket_thenReturnNoContent() {
        UUID id = UUID.randomUUID();
        Mockito.doNothing().when(ticketService).delete(id);

        mockMvc.perform(MockMvcRequestBuilders.delete("/api/v1/tickets/{id}", id))
                .andExpect(MockMvcResultMatchers.status().isNoContent());

        Mockito.verify(ticketService, Mockito.times(1)).delete(id);
    }

    @Test
    @WithMockUser(username = "johndoe@mail.com", roles = {"USER"})
    @SneakyThrows
    void givenValidCode_whenGetById_thenReturnResponse() {
        UUID id = UUID.randomUUID();
        UserDTO.UserResponse userResponse = UserDTO.UserResponse.fromUser(mockUser);
        TicketDTO.TicketResponse response = TicketDTO.TicketResponse.builder()
                .id(id.toString())
                .code("TCK-001")
                .title("Sample Ticket")
                .status("OPEN")
                .priority("HIGH")
                .createdAt(OffsetDateTime.now())
                .user(userResponse)
                .build();
        Mockito.when(ticketService.findById(id)).thenReturn(response);

        mockMvc.perform(MockMvcRequestBuilders.get("/api/v1/tickets/{id}", id)
                        .contentType(MediaType.APPLICATION_JSON))
                .andExpect(MockMvcResultMatchers.status().isOk())
                .andExpect(MockMvcResultMatchers.jsonPath("$.title").value(response.getTitle()));

        Mockito.verify(ticketService, Mockito.times(1)).findById(id);
    }

    @Test
    @WithMockUser(username = "johndoe@mail.com", roles = {"USER"})
    @SneakyThrows
    void givenValidCode_whenGetByCode_thenReturnResponse() {
        UserDTO.UserResponse userResponse = UserDTO.UserResponse.fromUser(mockUser);
        TicketDTO.TicketResponse response = TicketDTO.TicketResponse.builder()
                .id(UUID.randomUUID().toString())
                .code("TCK-001")
                .title("Sample Ticket")
                .status("OPEN")
                .priority("HIGH")
                .createdAt(OffsetDateTime.now())
                .user(userResponse)
                .build();
        Mockito.when(ticketService.findByCode("TCK-001")).thenReturn(response);

        mockMvc.perform(MockMvcRequestBuilders.get("/api/v1/tickets/code/TCK-001")
                        .contentType(MediaType.APPLICATION_JSON))
                .andExpect(MockMvcResultMatchers.status().isOk())
                .andExpect(MockMvcResultMatchers.jsonPath("$.title").value(response.getTitle()));

        Mockito.verify(ticketService, Mockito.times(1)).findByCode("TCK-001");
    }

    @Test
    @WithMockUser(username = "johndoe@mail.com", roles = {"USER"})
    @SneakyThrows
    void givenValidRequest_whenSearchTicket_returnResponses() {
        UserDTO.UserResponse userResponse = UserDTO.UserResponse.fromUser(mockUser);
        TicketDTO.TicketResponse response = TicketDTO.TicketResponse.builder()
                .id(UUID.randomUUID().toString())
                .code("TCK-001")
                .title("Sample Ticket")
                .status("OPEN")
                .priority("HIGH")
                .createdAt(OffsetDateTime.now())
                .user(userResponse)
                .build();
        Mockito.when(ticketService.searchTicket(null, null, null, null)).thenReturn(List.of(response));

        mockMvc.perform(MockMvcRequestBuilders.get("/api/v1/tickets")
                        .contentType(MediaType.APPLICATION_JSON))
                .andExpect(MockMvcResultMatchers.status().isOk())
                .andExpect(MockMvcResultMatchers.jsonPath("$[0].title").value(response.getTitle()));

        Mockito.verify(ticketService, Mockito.times(1)).searchTicket(null, null, null, null);
    }

    @Test
    @WithMockUser(username = "janedoe@mail.com", roles = {"USER"})
    @SneakyThrows
    void givenValidRequest_whenGetMyTickets_thenReturnUserTickets() {
        UserDTO.UserResponse userResponse = UserDTO.UserResponse.fromUser(mockUser);
        TicketDTO.TicketResponse response = TicketDTO.TicketResponse.builder()
                .id(UUID.randomUUID().toString())
                .code("TCK-002")
                .title("My Ticket")
                .status("RESOLVED")
                .priority("LOW")
                .createdAt(OffsetDateTime.now())
                .user(userResponse)
                .build();
        Mockito.when(ticketService.searchMyTicket(null, null, null, null)).thenReturn(List.of(response));

        mockMvc.perform(MockMvcRequestBuilders.get("/api/v1/tickets/me")
                        .contentType(MediaType.APPLICATION_JSON))
                .andExpect(MockMvcResultMatchers.status().isOk())
                .andExpect(MockMvcResultMatchers.jsonPath("$[0].title").value(response.getTitle()))
                .andExpect(MockMvcResultMatchers.jsonPath("$[0].code").value(response.getCode()));

        Mockito.verify(ticketService, Mockito.times(1)).searchMyTicket(null, null, null, null);
    }

    @Test
    @WithMockUser(username = "johndoe@mail.com", roles = {"USER"})
    @SneakyThrows
    void givenValidRequest_whenCreateReply_thenReturnReply() throws Exception {
        String code = "TCK-006";
        TicketCommentDTO.TicketCommentRequest request = TicketCommentDTO.TicketCommentRequest.builder()
                .content("Reply Message")
                .build();
        TicketCommentDTO.TicketCommentResponse response = TicketCommentDTO.TicketCommentResponse.builder()
                .id(UUID.randomUUID().toString())
                .content("Reply Message")
                .createdAt(OffsetDateTime.now())
                .build();
        Mockito.when(ticketCommentService.create(Mockito.eq(code), Mockito.any(TicketCommentDTO.TicketCommentRequest.class))).thenReturn(response);

        mockMvc.perform(MockMvcRequestBuilders.post("/api/v1/tickets/{code}/comments", code)
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(objectMapper.writeValueAsString(request)))
                .andExpect(MockMvcResultMatchers.status().isOk());

        Mockito.verify(ticketCommentService, Mockito.times(1)).create(Mockito.eq(code), Mockito.any(TicketCommentDTO.TicketCommentRequest.class));
    }
}