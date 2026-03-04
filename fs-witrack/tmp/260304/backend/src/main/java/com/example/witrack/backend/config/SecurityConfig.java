package com.example.witrack.backend.config;

import com.example.witrack.backend.security.AuthenticationEntryPointImpl;
import com.example.witrack.backend.security.AuthenticationTokenFilter;
import com.example.witrack.backend.security.jwt.JwtProvider;
import com.example.witrack.backend.security.user.UserDetailsServiceImpl;
import lombok.RequiredArgsConstructor;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.http.HttpMethod;
import org.springframework.security.authentication.AuthenticationManager;
import org.springframework.security.config.annotation.authentication.configuration.AuthenticationConfiguration;
import org.springframework.security.config.annotation.method.configuration.EnableMethodSecurity;
import org.springframework.security.config.annotation.web.builders.HttpSecurity;
import org.springframework.security.crypto.bcrypt.BCryptPasswordEncoder;
import org.springframework.security.crypto.password.PasswordEncoder;
import org.springframework.security.web.SecurityFilterChain;
import org.springframework.security.web.authentication.UsernamePasswordAuthenticationFilter;

@Configuration
@EnableMethodSecurity(prePostEnabled = true)
@RequiredArgsConstructor
public class SecurityConfig {

    private final AuthenticationEntryPointImpl authenticationEntryPoint;
    private final UserDetailsServiceImpl userDetailsService;
    private final JwtProvider jwtProvider;

    @Bean
    public PasswordEncoder passwordEncoder() {
        return new BCryptPasswordEncoder();
    }

    @Bean
    public AuthenticationManager authenticationManager(final AuthenticationConfiguration authenticationConfiguration) throws Exception {
        return authenticationConfiguration.getAuthenticationManager();
    }

    @Bean
    public AuthenticationTokenFilter authenticationTokenFilter() {
        return new AuthenticationTokenFilter(jwtProvider, userDetailsService);
    }

    @Bean
    public SecurityFilterChain configure(final HttpSecurity http) throws Exception {
        return http
                .cors(cors -> {
                }) // can be enabled if needed
                .csrf(csrf -> csrf.disable()) // must be disabled for a stateless API
                .authorizeHttpRequests(auth -> auth
                        .requestMatchers(
                                "/swagger-ui.html",
                                "/swagger-ui/**",
                                "/v3/api-docs/**"
                        ).permitAll()
                        .requestMatchers("/api/v1/auth/sign-up").permitAll()
                        .requestMatchers("/api/v1/auth/sign-in").permitAll()
                        .requestMatchers(HttpMethod.DELETE, "/api/v1/tickets/**").hasRole("ADMIN")
                        .requestMatchers(HttpMethod.PUT, "/api/v1/tickets/*/status").hasRole("ADMIN")
                        .anyRequest().authenticated()
                )
                // filter jwt only for api
                .addFilterBefore(authenticationTokenFilter(), UsernamePasswordAuthenticationFilter.class)
                .exceptionHandling(customizer -> customizer.authenticationEntryPoint(authenticationEntryPoint))
                .build();
    }
}
