

extern crate libc;
use std::ffi::CStr;





#[no_mangle]
pub extern "C" fn run_job(p_job_name: *const libc::c_char) {
    



    let buf_job_name = unsafe { 
        CStr::from_ptr(p_job_name).to_bytes() 
    };
    let job_name_str = String::from_utf8(buf_job_name.to_vec()).unwrap();
    
    

    println!("RUST ---------- running job - {}", job_name_str);
}
