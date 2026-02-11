package com.example.bepusdt.service;

import com.example.bepusdt.dto.PaymentNotifyRequest;
import com.example.bepusdt.util.SignUtil;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Service;

import java.util.HashMap;
import java.util.Map;

@Service
public class PaymentNotifyService {

    @Value("${bepusdt.api-token}")
    private String apiToken;

    /**
     * 回调签名验证
     */
    public boolean verifySignature(PaymentNotifyRequest req) {
        Map<String, Object> params = new HashMap<>();

        params.put("trade_id", req.getTrade_id());
        params.put("order_id", req.getOrder_id());
        params.put("amount", req.getAmount());
        params.put("actual_amount", req.getActual_amount());
        params.put("token", req.getToken());
        params.put("block_transaction_id", req.getBlock_transaction_id());
        params.put("status", req.getStatus());

        String localSign = SignUtil.sign(params, apiToken);

        return localSign.equalsIgnoreCase(req.getSignature());
    }
}
