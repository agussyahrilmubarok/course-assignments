package com.example.catalog.service;

import com.example.catalog.model.ProductDTO;

import java.util.List;

public interface CatalogService {

    ProductDTO.Response registerProduct(ProductDTO.RegisterRequest param);

    ProductDTO.Response findProductById(String productId);

    List<ProductDTO.Response> findProductsBySellerId(String sellerId);

    void deleteProduct(String productId);

    ProductDTO.Response decreaseStockCount(String productId, Long decreaseCount);
}
