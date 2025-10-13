package com.example.catalog.cassandra.domain;

import lombok.Getter;
import lombok.Setter;
import org.springframework.data.cassandra.core.mapping.Column;
import org.springframework.data.cassandra.core.mapping.PrimaryKey;
import org.springframework.data.cassandra.core.mapping.Table;

import java.util.List;

@Table("products")
@Getter
@Setter
public class Product {

    @PrimaryKey
    private String id;

    @Column("seller_id")
    private String sellerId;

    @Column("name")
    private String name;

    @Column("description")
    private String description;

    @Column("price")
    private Long price;

    @Column("stock_count")
    private Long stockCount;

    @Column("tags")
    private List<String> tags;

    public Product(String id, String sellerId, String name, String description, Long price, Long stockCount, List<String> tags) {
        this.id = id;
        this.sellerId = sellerId;
        this.name = name;
        this.description = description;
        this.price = price;
        this.stockCount = stockCount;
        this.tags = tags;
    }

    public Product() {}
}
