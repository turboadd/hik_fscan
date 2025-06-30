cmake_minimum_required(VERSION 3.10)
project(hik_adapter_static)

# ⚙️ บังคับใช้ C++11 ขึ้นไป
set(CMAKE_CXX_STANDARD 11)
set(CMAKE_CXX_STANDARD_REQUIRED True)

# 🔍 เพิ่ม source file ได้ง่ายในอนาคต
file(GLOB HIK_SRC
    "${CMAKE_CURRENT_SOURCE_DIR}/src/*.cpp"
)

# 🔨 สร้าง static library (.a)
add_library(hik_adapter STATIC ${HIK_SRC})

# 📂 รวม path header (.h) ต่างๆ
target_include_directories(hik_adapter
    PUBLIC ${CMAKE_CURRENT_SOURCE_DIR}/includes
)

# 🔗 ลิงก์กับ HCNetSDK.lib
target_link_directories(hik_adapter
    PRIVATE ${CMAKE_CURRENT_SOURCE_DIR}/libs
)
target_link_libraries(hik_adapter
    PRIVATE HCNetSDK
)

# 📦 output .a ไปไว้ใน libs/ สำหรับ Go ใช้
set_target_properties(hik_adapter PROPERTIES
    ARCHIVE_OUTPUT_DIRECTORY ${CMAKE_CURRENT_SOURCE_DIR}/libs
)
