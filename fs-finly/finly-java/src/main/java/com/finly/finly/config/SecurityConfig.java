package com.finly.finly.config;

import jakarta.servlet.Filter;
import lombok.RequiredArgsConstructor;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.security.authentication.AuthenticationManager;
import org.springframework.security.config.annotation.authentication.configuration.AuthenticationConfiguration;
import org.springframework.security.config.annotation.web.builders.HttpSecurity;
import org.springframework.security.config.annotation.web.configuration.EnableWebSecurity;
import org.springframework.security.config.annotation.web.configurers.AbstractHttpConfigurer;
import org.springframework.security.config.http.SessionCreationPolicy;
import org.springframework.security.crypto.bcrypt.BCryptPasswordEncoder;
import org.springframework.security.crypto.password.PasswordEncoder;
import org.springframework.security.web.SecurityFilterChain;
import org.springframework.security.web.context.SecurityContextPersistenceFilter;

@Configuration
@EnableWebSecurity
@RequiredArgsConstructor
public class SecurityConfig {

    @Bean
    public SecurityFilterChain securityFilterChain(HttpSecurity httpSecurity) throws Exception {
        httpSecurity
                .csrf(AbstractHttpConfigurer::disable)
                .authorizeHttpRequests(auth -> auth
                        .requestMatchers("/login", "/signup").anonymous()
                        .requestMatchers("/", "/favicon.ico", "/css/**", "/js/**", "/images/**").permitAll()
                        .anyRequest().authenticated() // Require authentication for all other requests
                )
                .formLogin(form -> form.disable())
                .logout(out -> out
                        .logoutSuccessUrl("/login")
                        .invalidateHttpSession(true)
                        .deleteCookies("JSESSIONID")
                        .permitAll()
                )
                .sessionManagement(session -> session
                        .sessionCreationPolicy(SessionCreationPolicy.IF_REQUIRED)
                        .invalidSessionUrl("/login")
                )
                .exceptionHandling(exception -> exception
                        // Handle Access Denied by redirecting to login page (or any custom page you like)
                        .accessDeniedPage("/login")  // Redirects to login page instead of showing "Access Denied"
                );
        return httpSecurity.build();
    }

    @Bean
    public Filter securityContextPersistenceFilter() {
        // The SecurityContextPersistenceFilter is responsible for storing and restoring
        // the SecurityContext (authentication information) in the session during each request.
        // This is crucial for ensuring that the authentication context is persisted across multiple requests
        // and that the user remains authenticated while interacting with the application.
        return new SecurityContextPersistenceFilter();
    }

    @Bean
    public AuthenticationManager authenticationManager(AuthenticationConfiguration authenticationConfiguration) throws Exception {
        return authenticationConfiguration.getAuthenticationManager();
    }

    @Bean
    public PasswordEncoder passwordEncoder() {
        return new BCryptPasswordEncoder();
    }
}
