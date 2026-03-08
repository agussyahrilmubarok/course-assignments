package com.example.graphql.rest;

import com.example.graphql.model.BookDTO;
import com.example.graphql.service.BookService;
import io.swagger.v3.oas.annotations.responses.ApiResponse;
import jakarta.validation.Valid;
import org.springframework.hateoas.CollectionModel;
import org.springframework.hateoas.EntityModel;
import org.springframework.http.HttpStatus;
import org.springframework.http.MediaType;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.List;
import java.util.stream.Collectors;

import static org.springframework.hateoas.server.mvc.WebMvcLinkBuilder.linkTo;
import static org.springframework.hateoas.server.mvc.WebMvcLinkBuilder.methodOn;

@RestController
@RequestMapping(value = "/api/books", produces = MediaType.APPLICATION_JSON_VALUE)
public class BookResource {

    private final BookService bookService;

    public BookResource(BookService bookService) {
        this.bookService = bookService;
    }

    @GetMapping
    public ResponseEntity<CollectionModel<EntityModel<BookDTO.Response>>> getAllBooks() {
        List<EntityModel<BookDTO.Response>> books = bookService.findAll()
                .stream()
                .map(book -> EntityModel.of(book,
                        linkTo(methodOn(this.getClass()).getBook(book.getId())).withSelfRel(),
                        linkTo(methodOn(this.getClass()).getAllBooks()).withRel("all-books")
                ))
                .collect(Collectors.toList());

        CollectionModel<EntityModel<BookDTO.Response>> collectionModel = CollectionModel.of(books,
                linkTo(methodOn(this.getClass()).getAllBooks()).withSelfRel());

        return ResponseEntity.ok(collectionModel);
    }

    @GetMapping("/{id}")
    public ResponseEntity<EntityModel<BookDTO.Response>> getBook(@PathVariable(name = "id") final Long id) {
        BookDTO.Response bookDTO = bookService.get(id);

        EntityModel<BookDTO.Response> entityModel = EntityModel.of(bookDTO,
                linkTo(methodOn(this.getClass()).getBook(id)).withSelfRel(),
                linkTo(methodOn(this.getClass()).getAllBooks()).withRel("all-books")
        );

        return ResponseEntity.ok(entityModel);
    }

    @PostMapping
    @ApiResponse(responseCode = "201")
    public ResponseEntity<EntityModel<BookDTO.Response>> createBook(@RequestBody @Valid final BookDTO.BookRequest payload) {
        BookDTO.Response bookDTO = bookService.create(payload);

        EntityModel<BookDTO.Response> entityModel = EntityModel.of(bookDTO,
                linkTo(methodOn(this.getClass()).getBook(bookDTO.getId())).withSelfRel(),
                linkTo(methodOn(this.getClass()).getAllBooks()).withRel("all-books")
        );

        return new ResponseEntity<>(entityModel, HttpStatus.CREATED);
    }

    @PutMapping("/{id}")
    public ResponseEntity<EntityModel<BookDTO.Response>> updateBook(@PathVariable(name = "id") final Long id,
                                                                    @RequestBody @Valid final BookDTO.BookRequest payload) {
        BookDTO.Response bookDTO = bookService.update(id, payload);

        EntityModel<BookDTO.Response> entityModel = EntityModel.of(bookDTO,
                linkTo(methodOn(this.getClass()).getBook(id)).withSelfRel(),
                linkTo(methodOn(this.getClass()).getAllBooks()).withRel("all-books")
        );

        return ResponseEntity.ok(entityModel);
    }

    @DeleteMapping("/{id}")
    @ApiResponse(responseCode = "204")
    public ResponseEntity<Void> deleteBook(@PathVariable(name = "id") final Long id) {
        bookService.delete(id);

        return ResponseEntity.noContent().build();
    }
}
