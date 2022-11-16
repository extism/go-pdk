#pragma once

#include <stdint.h>

#define IMPORT(a, b) __attribute__((import_module(a), import_name(b)))

IMPORT("env", "extism_input_length") extern uint64_t extism_input_length();
IMPORT("env", "extism_length") extern uint64_t extism_length(uint64_t);
IMPORT("env", "extism_alloc") extern uint64_t extism_alloc(uint64_t);
IMPORT("env", "extism_free") extern void extism_free(uint64_t);

IMPORT("env", "extism_input_load_u8")
extern uint8_t extism_input_load_u8(uint64_t);

IMPORT("env", "extism_input_load_u64")
extern uint64_t extism_input_load_u64(uint64_t);

IMPORT("env", "extism_output_set")
extern void extism_output_set(uint64_t, uint64_t);

IMPORT("env", "extism_error_set")
extern void extism_error_set(uint64_t);

IMPORT("env", "extism_config_get")
extern uint64_t extism_config_get(uint64_t);

IMPORT("env", "extism_var_get")
extern uint64_t extism_var_get(uint64_t);

IMPORT("env", "extism_var_set")
extern void extism_var_set(uint64_t, uint64_t);

IMPORT("env", "extism_store_u8")
extern void extism_store_u8(uint64_t, uint8_t);

IMPORT("env", "extism_load_u8")
extern uint8_t extism_load_u8(uint64_t);

IMPORT("env", "extism_store_u64")
extern void extism_store_u64(uint64_t, uint64_t);

IMPORT("env", "extism_load_u64")
extern uint64_t extism_load_u64(uint64_t);

IMPORT("env", "extism_http_request")
extern uint64_t extism_http_request(uint64_t, uint64_t);

IMPORT("env", "extism_log_info")
extern void extism_log_info(uint64_t);
IMPORT("env", "extism_log_debug")
extern void extism_log_debug(uint64_t);
IMPORT("env", "extism_log_warn")
extern void extism_log_warn(uint64_t);
IMPORT("env", "extism_log_error")
extern void extism_log_error(uint64_t);

static void extism_load(uint64_t offs, uint8_t *buffer, uint64_t length) {
  uint64_t n;
  uint64_t left = 0;

  for (uint64_t i = 0; i < length; i += 1) {
    left = length - i;
    if (left < 8) {
      buffer[i] = extism_load_u8(offs + i);
      continue;
    }

    n = extism_load_u64(offs + i);
    *((uint64_t *)buffer + (i / 8)) = n;
    i += 7;
  }
}

static void extism_load_input(uint8_t *buffer, uint64_t length) {
  uint64_t n;
  uint64_t left = 0;

  for (uint64_t i = 0; i < length; i += 1) {
    left = length - i;
    if (left < 8) {
      buffer[i] = extism_input_load_u8(i);
      continue;
    }

    n = extism_input_load_u64(i);
    *((uint64_t *)buffer + (i / 8)) = n;
    i += 7;
  }
}

static void extism_store(uint64_t offs, const uint8_t *buffer,
                         uint64_t length) {
  uint64_t n;
  uint64_t left = 0;
  for (uint64_t i = 0; i < length; i++) {
    left = length - i;
    if (left < 8) {
      extism_store_u8(offs + i, buffer[i]);
      continue;
    }

    n = *((uint64_t *)buffer + (i / 8));
    extism_store_u64(offs + i, n);
    i += 7;
  }
}
