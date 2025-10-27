package com.example.graphql.controller;

import com.example.graphql.model.AuthorDTO;
import com.example.graphql.service.AuthorService;
import lombok.RequiredArgsConstructor;
import org.springframework.graphql.data.method.annotation.Argument;
import org.springframework.graphql.data.method.annotation.MutationMapping;
import org.springframework.graphql.data.method.annotation.QueryMapping;
import org.springframework.stereotype.Controller;

import java.util.List;

@Controller
@RequiredArgsConstructor
public class AuthorController {

    private final AuthorService authorService;

    @QueryMapping
    public List<AuthorDTO.Response> findAllAuthors() {
        return authorService.findAll();
    }

    @QueryMapping
    public AuthorDTO.Response getAuthor(@Argument Long id) {
        return authorService.get(id);
    }

    @MutationMapping
    public AuthorDTO.Response createAuthor(@Argument("input") AuthorDTO.AuthorRequest input) {
        return authorService.create(input);
    }

    @MutationMapping
    public AuthorDTO.Response updateAuthor(@Argument Long id, @Argument("input") AuthorDTO.AuthorRequest input) {
        return authorService.update(id, input);
    }

    @MutationMapping
    public Boolean deleteAuthor(@Argument Long id) {
        authorService.delete(id);
        return true;
    }
}
