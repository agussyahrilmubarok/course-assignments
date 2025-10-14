package com.example.catalog.service;

import com.example.catalog.cassandra.domain.Product;
import com.example.catalog.cassandra.repos.ProductRepository;
import com.example.catalog.exception.InsufficientStockException;
import com.example.catalog.exception.ProductNotFoundException;
import com.example.catalog.model.ProductDTO;
import com.example.catalog.postgres.domain.SellerProduct;
import com.example.catalog.postgres.repos.SellerProductRepository;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.util.List;
import java.util.UUID;

@Service("CatalogServiceImpl")
@Slf4j
@RequiredArgsConstructor
public class CatalogServiceImpl implements CatalogService {

    private final ProductRepository productRepository;
    private final SellerProductRepository sellerProductRepository;

    @Override
    @Transactional
    public ProductDTO.Response registerProduct(ProductDTO.RegisterRequest param) {
        String productId = UUID.randomUUID().toString();
        String sellerProductId = UUID.randomUUID().toString();

        SellerProduct sellerProduct = new SellerProduct();
        sellerProduct.setId(sellerProductId);
        sellerProduct.setSellerId(param.getSellerId());
        sellerProduct.setProductId(productId);
        sellerProductRepository.save(sellerProduct);
        log.info("Saved seller product with ID: {}", sellerProductId);

        Product product = new Product();
        product.setId(productId);
        product.setSellerId(param.getSellerId());
        product.setName(param.getName());
        product.setDescription(param.getDescription());
        product.setPrice(param.getPrice());
        product.setStockCount(param.getStockCount());
        product.setTags(param.getTags());
        product = productRepository.save(product);
        log.info("Saved product with ID: {}", productId);

        return ProductDTO.Response.from(product);
    }

    @Override
    public ProductDTO.Response findProductById(String productId) {
        log.info("Fetching product with ID: {}", productId);
        Product product = productRepository.findById(productId)
                .orElseThrow(() -> {
                    log.error("Product with ID {} not found", productId);
                    return new ProductNotFoundException("Product not found");
                });

        return ProductDTO.Response.from(product);
    }

    @Override
    public List<ProductDTO.Response> findProductsBySellerId(String sellerId) {
        log.info("Fetching products for sellerId: {}", sellerId);
        return productRepository.findAllBySellerId(sellerId)
                .stream()
                .map(ProductDTO.Response::from)
                .toList();
    }

    @Override
    @Transactional
    public void deleteProduct(String productId) {
        productRepository.deleteById(productId);
        sellerProductRepository.deleteByProductId(productId);
        log.info("Deleted product and associated sellerProduct for ID: {}", productId);
    }

    @Override
    public ProductDTO.Response decreaseStockCount(String productId, Long decreaseCount) {
        Product product = productRepository.findById(productId)
                .orElseThrow(() -> {
                    log.error("Product with ID {} not found for stock decrease", productId);
                    return new ProductNotFoundException("Product not found");
                });

        if (product.getStockCount() < decreaseCount) {
            log.error("Insufficient stock for productId: {}. Current: {}, Requested: {}",
                    productId, product.getStockCount(), decreaseCount);
            throw new InsufficientStockException("Insufficient stock");
        }

        product.setStockCount(product.getStockCount() - decreaseCount);
        product = productRepository.save(product);
        log.info("Stock updated for productId: {}. New stock: {}", productId, product.getStockCount());

        return ProductDTO.Response.from(product);
    }
}
