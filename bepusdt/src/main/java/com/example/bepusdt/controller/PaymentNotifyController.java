package com.example.bepusdt.controller;

import com.example.bepusdt.dto.PaymentNotifyRequest;
import com.example.bepusdt.service.PaymentNotifyService;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.web.bind.annotation.*;

@RestController
@RequestMapping("/api/payment")
public class PaymentNotifyController {

    private static final Logger log = LoggerFactory.getLogger(PaymentNotifyController.class);

    private final PaymentNotifyService notifyService;

    public PaymentNotifyController(PaymentNotifyService notifyService) {
        this.notifyService = notifyService;
    }

    /**
     * BEpusdt 支付回调
     */
    @PostMapping("/notify")
    public String notify(@RequestBody PaymentNotifyRequest request) {

        // 1️⃣ 基础日志（必须）
        log.info("BEpusdt notify received, orderId={}, status={}",
                request.getOrder_id(), request.getStatus());

        // 2️⃣ 验签
        boolean verified = notifyService.verifySignature(request);
        if (!verified) {
            log.error("BEpusdt notify signature verify FAILED, orderId={}",
                    request.getOrder_id());

            // ⚠️ 仍然返回 ok，防止对方无限重试
            return "ok";
        }

        // 3️⃣ 根据状态分支处理（后续你会在这里接 DB / MQ）
        switch (request.getStatus()) {
            case 1:
                log.info("Order waiting payment: {}", request.getOrder_id());
                break;

            case 2:
                log.info("Order payment SUCCESS: {}", request.getOrder_id());
                // TODO: 幂等校验 + 入账
                break;

            case 3:
                log.info("Order payment TIMEOUT: {}", request.getOrder_id());
                break;

            default:
                log.warn("Unknown payment status: {}", request.getStatus());
        }

        // 4️⃣ 按 BEpusdt 规范，成功必须返回 ok
        return "ok";
    }
}
