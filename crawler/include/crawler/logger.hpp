/*
 * Copyright (c) 2013-2016 John Connor
 * Copyright (c) 2016-2017 The Pila Developers
 *
 * This file is part of Pila.
 *
 * Pila is free software: you can redistribute it and/or modify
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
 
#ifndef CRAWLER_LOGGER_HPP
#define CRAWLER_LOGGER_HPP

#if (defined __ANDROID__)
#include <android/log.h>
#endif
  
#if (defined _WIN32 || defined WIN32) || (defined _WIN64 || defined WIN64)
#define WIN32_LEAN_AND_MEAN
#include <windows.h>
#endif // (defined _WIN32 || defined WIN32) || (defined _WIN64 || defined WIN64)

#include <iostream>
#include <fstream>
#include <memory>
#include <sstream>

#include <mutex>

namespace crawler {

    /**
     * Implements a logger.
     */
    class logger
    {
        public:
            
			typedef enum severity
			{
				severity_none,
				severity_debug,
				severity_error,
				severity_info,
				severity_warning,
                severity_test,
			} severity_t;
			
            /**
             * Singleton accessor.
             */
			static logger & instance()
			{
			    static logger g_logger;
                
			    return g_logger;
			}
            
            /**
             * operator <<
             */
            template <class T>
            logger & operator << (T const & val)
            {
                std::stringstream ss;
                
                ss << val;
                
                log(ss);
                
                ss.str(std::string());
                
                return logger::instance();
            }

            /**
             * Perform the actual logging.
             * @param val
             */
			void log(std::stringstream & val)
			{
			    static bool use_cout = true;

			    if (use_cout)
			    {
#if (defined _WIN32 || defined WIN32) || (defined _WIN64 || defined WIN64)
#if defined(_UNICODE)
			        DWORD len = MultiByteToWideChar(
			            CP_ACP, 0, val.str().c_str(), -1, NULL, 0
			        );

			        std::unique_ptr<wchar_t> buf(new wchar_t[len]);

			        MultiByteToWideChar(
			            CP_ACP, 0, val.str().c_str(), -1, buf.get(), len
			        );

			        OutputDebugString(buf.get());
			        OutputDebugString(L"\n");

			        std::cerr << val.str() << std::endl;
#else
			        OutputDebugString(val.str().c_str());
			        OutputDebugString(L"\n");

			        std::cerr << val.str() << std::endl;
#endif // _UNICODE
#else // Not Windows.
#if (defined __ANDROID__)
					__android_log_print(
                        ANDROID_LOG_DEBUG, "logger", val.str().c_str()
                    );
#else
			        std::cerr << val.str() << std::endl;
#endif
#endif // defined _WIN32 || defined WIN32) || (defined _WIN64 || defined WIN64
			    }
			}
            
        private:
        
			// ...
            
        protected:
        
            /**
             * The std::ofstream.
             */
            std::ofstream ofstream_;
    };
    
    #define log_xx(severity, strm) \
    { \
        std::stringstream __ss; \
        switch (severity) \
        { \
            case crawler::logger::severity_debug: \
                __ss << "[DEBUG] - "; \
            break; \
            case crawler::logger::severity_error: \
                __ss << "[ERROR] - "; \
            break; \
            case crawler::logger::severity_info: \
                __ss << "[INFO] - "; \
            break; \
            case crawler::logger::severity_warning: \
                __ss << "[WARNING] - "; \
            break; \
            case crawler::logger::severity_test: \
                __ss << "[TEST] - "; \
            break; \
            default: \
                __ss << "[UNKNOWN] - "; \
        } \
		__ss << __FUNCTION__ << ": "; \
        __ss << strm; \
        crawler::logger::instance() << __ss.str(); \
        __ss.str(std::string()); \
    } \
	
#define log_none(strm) /** */
#if (defined NDEBUG)
#define log_debug(strm) log_none(strm)
#else
#define log_debug(strm) log_xx(crawler::logger::severity_debug, strm)
#endif
#define log_error(strm) log_xx(crawler::logger::severity_error, strm)
#define log_info(strm) log_xx(crawler::logger::severity_info, strm)
#define log_warn(strm) log_xx(crawler::logger::severity_warning, strm)
#define log_test(strm) log_xx(crawler::logger::severity_test, strm)

} // namespace crawler

#endif // CRAWLER_LOGGER_HPP

        