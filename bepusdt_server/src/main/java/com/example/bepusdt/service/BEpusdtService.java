package com.example.bepusdt.service;

import com.example.bepusdt.dto.CancelOrderRequest;
import com.example.bepusdt.dto.CreateOrderRequest;
import com.example.bepusdt.util.SignUtil;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.http.*;
import org.springframework.stereotype.Service;
import org.springframework.web.client.RestTemplate;
//import tools.jackson.databind.ObjectMapper;

import java.util.HashMap;
import java.util.Map;

@Service
public class BEpusdtService {

    @Value("${bepusdt.base-url}")
    private String baseUrl;

    @Value("${bepusdt.api-token}")
    private String apiToken;

    private final RestTemplate restTemplate = new RestTemplate();
//    private final ObjectMapper objectMapper = new ObjectMapper();

    public String createOrder(CreateOrderRequest req) {
        Map<String, Object> params = toCreateOrderMap(req);
        params.put("signature", SignUtil.sign(params, apiToken));
        return post("/api/v1/order/create-transaction", params);
    }

    public String cancelOrder(CancelOrderRequest req) {
        Map<String, Object> params = new HashMap<>();
        params.put("trade_id", req.getTradeId());
        params.put("signature", SignUtil.sign(params, apiToken));
        return post("/api/v1/order/cancel-transaction", params);
    }

    private String post(String path, Object body) {
        HttpHeaders headers = new HttpHeaders();
        headers.setContentType(MediaType.APPLICATION_JSON);
        return restTemplate.postForObject(
                baseUrl + path,
                new HttpEntity<>(body, headers),
                String.class
        );
    }

    private Map<String, Object> toCreateOrderMap(CreateOrderRequest r) {
        Map<String, Object> m = new HashMap<>();
        m.put("order_id", r.getOrderId());
        m.put("amount", r.getAmount());
        m.put("notify_url", r.getNotifyUrl());
        m.put("redirect_url", r.getRedirectUrl());
        m.put("trade_type", r.getTradeType());
        m.put("fiat", r.getFiat());
        m.put("address", r.getAddress());
        m.put("name", r.getName());
        m.put("timeout", r.getTimeout());
        m.put("rate", r.getRate());
        return m;
    }
}


