/*
 * Copyright (c) 2013-2016 John Connor
 * Copyright (c) 2016-2017 The Pila developers
 *
 * This file is part of pila.
 *
 * pila is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License with
 * additional permissions to the one published by the Free Software
 * Foundation, either version 3 of the License, or (at your option)
 * any later version. For more information see LICENSE.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program. If not, see <http://www.gnu.org/licenses/>.
 */

#ifndef COIN_ZEROTIME_VOTE_HPP
#define COIN_ZEROTIME_VOTE_HPP

#include <cstdint>
#include <string>
#include <vector>

#include <coin/data_buffer.hpp>
#include <coin/key_public.hpp>
#include <coin/sha256.hpp>
#include <coin/transaction_in.hpp>

namespace coin {

    /**
     * Implements a ZeroTime vote message (ztvote).
     */
    class zerotime_vote : public data_buffer
    {
        public:
        
            /**
             * Constructor
             */
            zerotime_vote();
        
            /**
             * Constructor
             * @param block_height The block height.
             * @param hash_block The block hash.
             * @param hash_tx The transaction hash.
             * @param transactions_in The transaction_in's.
             * @param public_key The key_public.
             */
            zerotime_vote(
                const std::uint32_t & block_height, const sha256 & hash_block,
                const sha256 & hash_tx,
                const std::vector<transaction_in> & transactions_in,
                const key_public & public_key
            );
        
            /**
             * Encodes
             */
            void encode();
        
            /**
             * Encodes
             * @param buffer The data_buffer.
             * @param is_copy If true we are encoding a copy.
             */
            void encode(data_buffer & buffer, const bool & is_copy = false);
        
            /**
             * Decodes
             */
            bool decode();
        
            /**
             * Encodes
             * @param buffer The data_buffer.
             */
            bool decode(data_buffer & buffer);
        
            /**
             * Set's null.
             */
            void set_null();
        
            /**
             * The (current) block height.
             */
            const std::uint32_t & block_height() const;
        
            /**
             * The hash of the (current) block (liar prevention).
             */
            const sha256 & hash_block() const;
        
            /**
             * The hash of the transaction we are voting for.
             */
            const sha256 & hash_tx() const;
        
            /**
             * The hash nonce (random).
             */
            const sha256 & hash_nonce() const;
        
            /**
             * The inputs we are voting for.
             */
            const std::vector<transaction_in> & transactions_in() const;
        
            /**
             * The public key.
             */
            const key_public & public_key() const;
        
            /**
             * The (calculated) score.
             */
            const std::int16_t & score() const;
        
        private:

            /**
             * Signs
             * @param buffer The data_buffer.
             */
            bool sign(data_buffer & buffer);
        
            /**
             * Verifies
             * @param buffer The data_buffer.
             */
            bool verify(data_buffer & buffer);
        
            /**
             * The version.
             */
            enum { current_version = 1 };

            /**
             * The version.
             */
            std::uint32_t m_version;
        
            /**
             * The (current) block height.
             */
            std::uint32_t m_block_height;
        
            /**
             * The hash of the (current) block (liar prevention).
             */
            sha256 m_hash_block;
        
            /**
             * The hash of the transaction we are voting for.
             */
            sha256 m_hash_tx;
        
            /**
             * The hash nonce (random).
             */
            sha256 m_hash_nonce;
        
            /**
             * The inputs we are voting for.
             */
            std::vector<transaction_in> m_transactions_in;

            /**
             * The public key.
             */
            key_public m_public_key;
        
            /**
             * The signature (calculated by the wallet address).
             */
            std::vector<std::uint8_t> m_signature;
        
            /**
             * The (calculated) score.
             */
            mutable std::int16_t m_score;

        protected:
        
            // ...
    };
    
} // namespace coin

#endif // COIN_ZEROTIME_VOTE_HPP
