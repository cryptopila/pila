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

#ifndef COIN_RPC_TRANSPORT_HPP
#define COIN_RPC_TRANSPORT_HPP

#include <chrono>
#include <cstdint>
#include <deque>
#include <memory>

#if (defined __IPHONE_OS_VERSION_MAX_ALLOWED)
#import <CFNetwork/CFSocketStream.h>
#endif // __IPHONE_OS_VERSION_MAX_ALLOWED

#include <boost/asio.hpp>

namespace coin {

    /**
     * Implements a tcp transport.
     */
    class rpc_transport
        : public std::enable_shared_from_this<rpc_transport>
    {
        public:
        
            /**
             * The states.
             */
            typedef enum
            {
                state_disconnected,
                state_connecting,
                state_connected,
            } state_t;
        
            /**
             * Constructor
             * @param ios The boost::asio::io_service.
             * @param s The boost::asio::strand.
             */
            rpc_transport(
                boost::asio::io_service & ios, boost::asio::strand & s
            );

            /**
             * Destructor
             */
            ~rpc_transport();

            /**
             * Starts the transport (outgoing).
             * f The completion handler.
             */
            void start(
                const std::string & hostname, const std::uint16_t & port,
                const std::function<void (boost::system::error_code,
                std::shared_ptr<rpc_transport>)> &
            );
        
            /**
             * Starts the transport (incoming).
             */
            void start();
        
            /**
             * Stops the transport.
             */
            void stop();
        
            /**
             * Sets the on read handler.
             * @param f the std::function.
             */
            void set_on_read(
                const std::function<void (std::shared_ptr<rpc_transport>,
                const char *, const std::size_t &)> &
            );
        
            /**
             * Performs a write operation.
             * @param buf The buffer.
             * @param len The length.
             */
            void write(const char *, const std::size_t &);
        
            /**
             * The state.
             */
            state_t & state();
        
            /**
             * Sets the identifier.
             * @param val The value.
             */
            void set_identifier(const std::string &);
        
            /**
             * The identifier.
             */
            const std::string & identifier() const;
        
            /**
             * The socket.
             */
            boost::asio::ip::tcp::socket & socket();

            /**
             * If true the conneciton will close as soon as it's write queue is
             * exhausted.
             * @param flag The flag.
             */
            void set_close_after_writes(const bool &);
        
            /**
             * Sets the read timeout.
             * @param val The value.
             */
            void set_read_timeout(const std::uint32_t &);
        
            /**
             * Sets the write timeout.
             * @param val The value.
             */
            void set_write_timeout(const std::uint32_t &);
        
            /**
             * Runs the test case.
             */
            static int run_test();
        
        private:
        
            /**
             * do_connect
             */
            void do_connect(const boost::asio::ip::tcp::endpoint & ep);
        
            /**
             * do_connect
             */
            void do_connect(boost::asio::ip::tcp::resolver::iterator);
        
            /**
             * do_read
             */
            void do_read();
        
            /**
             * do_write
             */
            void do_write(const char * buf, const std::size_t & len);
        
            /**
             * The identifier.
             */
            std::string m_identifier;
        
            /**
             * The state.
             */
            state_t m_state;

            /**
             * The boost::asio::ip::tcp::socket.
             */
            std::shared_ptr<boost::asio::ip::tcp::socket> m_socket;

            /**
             * If true the conneciton will close as soon as it's write queue is
             * exhausted.
             */
            bool m_close_after_writes;
        
            /**
             * The read timeout.
             */
            std::uint32_t m_read_timeout;
        
            /**
             * The write timeout.
             */
            std::uint32_t m_write_timeout;
        
            /**
             * The completion handler.
             */
            std::function<
                void (boost::system::error_code, std::shared_ptr<rpc_transport>)
            > m_on_complete;
        
            /**
             * The read handler.
             */
            std::function<
                void (std::shared_ptr<rpc_transport>, const char *,
                const std::size_t &)
            > m_on_read;
            
        protected:
        
            friend class rpc_connection;
            
            /**
             * Set's up the socket for VoIP operation on iOS.
             */
            void set_voip();
        
            /**
             * The boost::asio::io_service.
             */
            boost::asio::io_service & io_service_;
        
            /**
             * The boost::asio::strand.
             */
            boost::asio::strand strand_;

            /**
             * The connect timeout timer.
             */
            boost::asio::basic_waitable_timer<
                std::chrono::steady_clock
            > connect_timeout_timer_;
        
            /**
             * The read timeout timer.
             */
            boost::asio::basic_waitable_timer<
                std::chrono::steady_clock
            > read_timeout_timer_;
        
            /**
             * The write timeout timer.
             */
            boost::asio::basic_waitable_timer<
                std::chrono::steady_clock
            > write_timeout_timer_;
        
            /**
             * The write queue.
             */
            std::deque< std::vector<char> > write_queue_;
        
            /**
             * The read buffer.
             */
            char read_buffer_[1024];
        
#if (defined __IPHONE_OS_VERSION_MAX_ALLOWED)
            /**
             * Read and Write streams for background voip usage on iOS.
             */
            CFReadStreamRef readStreamRef_;
            CFWriteStreamRef writeStreamRef_;
#endif // __IPHONE_OS_VERSION_MAX_ALLOWED
    };

} // namespace coin

#endif // COIN_RPC_TRANSPORT_HPP
