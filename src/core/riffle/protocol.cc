
#include "core/riffle/protocol.h"

#include <algorithm>  // std::random_shuffle
#include <vector>     // std::vector

#include "core/crypto/rand.h"

namespace kovri {
namespace core {

RiffleProtocol::RiffleProtocol() : elgamal_private_key_(new std::uint8_t[256]), elgamal_public_key_(new std::uint8_t[256]), signing_private_key_(new std::uint8_t[64]), signing_public_key_(new std::uint8_t[32]) {
  // Generate ElGamal and EdDSA key pairs
  kovri::core::GenerateElGamalKeyPair(elgamal_private_key_.get(), elgamal_public_key_.get());
  kovri::core::CreateEd25519KeyPair(signing_private_key_.get(), signing_public_key_.get());
}

RiffleProtocol::~RiffleProtocol() {}

void RiffleProtocol::Shuffle(
    const std::vector<std::vector<std::uint8_t>>& messages,
    std::vector<std::vector<std::uint8_t>>* shuffled_messages) {
  // Create a copy of the messages
  *shuffled_messages = messages;

  // Shuffle the messages
  kovri::core::Shuffle(shuffled_messages->begin(), shuffled_messages->end());
}

}  // namespace core
}  // namespace kovri
