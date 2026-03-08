package com.example.graphql.controller;

import com.example.graphql.model.BookDTO;
import com.example.graphql.service.BookService;
import lombok.RequiredArgsConstructor;
import org.springframework.graphql.data.method.annotation.*;
import org.springframework.stereotype.Controller;
import java.util.List;

@Controller
@RequiredArgsConstructor
public class BookController {

    private final BookService bookService;

    @QueryMapping
    public List<BookDTO.Response> findAllBooks() {
        return bookService.findAll();
    }

    @QueryMapping
    public BookDTO.Response getBook(@Argument Long id) {
        return bookService.get(id);
    }

    @MutationMapping
    public BookDTO.Response createBook(@Argument("input") BookDTO.BookRequest input) {
        return bookService.create(input);
    }

    @MutationMapping
    public BookDTO.Response updateBook(@Argument Long id, @Argument("input") BookDTO.BookRequest input) {
        return bookService.update(id, input);
    }

    @MutationMapping
    public Boolean deleteBook(@Argument Long id) {
        bookService.delete(id);
        return true;
    }
}
