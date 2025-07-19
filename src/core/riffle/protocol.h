
#ifndef SRC_CORE_RIFFLE_PROTOCOL_H_
#define SRC_CORE_RIFFLE_PROTOCOL_H_

#include <cstdint>
#include <vector>

#include "core/crypto/elgamal.h"
#include "core/crypto/signature.h"

namespace kovri {
namespace core {

class RiffleProtocol {
 public:
  RiffleProtocol();
  ~RiffleProtocol();

  // Verifiable shuffle
  void Shuffle(
      const std::vector<std::vector<std::uint8_t>>& messages,
      std::vector<std::vector<std::uint8_t>>* shuffled_messages);

 private:
  // ElGamal encryption key pair
  std::unique_ptr<std::uint8_t[]> elgamal_private_key_;
  std::unique_ptr<std::uint8_t[]> elgamal_public_key_;

  // EdDSA signing key pair
  std::unique_ptr<std::uint8_t[]> signing_private_key_;
  std::unique_ptr<std::uint8_t[]> signing_public_key_;
};

}  // namespace core
}  // namespace kovri

#endif  // SRC_CORE_RIFFLE_PROTOCOL_H_
