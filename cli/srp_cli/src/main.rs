extern crate num_bigint;
extern crate num_traits;
extern crate num_integer;
extern crate base64;
extern crate rand;
extern crate hkdf;

use std::io;
use std::io::Write;
use rand::rngs::OsRng;
use num_bigint::{BigInt, ToBigInt, RandBigInt, Sign};
use num_traits::sign::Signed;
use num_integer::Integer;
use sha2::Sha256;
use hmac::{Hmac, Mac};
use hkdf::Hkdf;

type HmacSha256 = Hmac<Sha256>;

#[allow(non_snake_case)]
trait PositiveModPow {
    fn pos_mod_pow(&self, exp: &BigInt, N: &BigInt) -> BigInt;
}

#[allow(non_snake_case)]
impl PositiveModPow for BigInt {
    fn pos_mod_pow(&self, exp: &BigInt, N: &BigInt) -> BigInt {
        let x: BigInt = 
            if self.is_negative() {
                self.clone()
            } else {
                self.mod_floor(&N)
            };

        return x.modpow(&exp, &N);
    }
}


static _N_BYTES: &[u8] = &[
    0xBD, 0xE5, 0xC8, 0x29, 0xE8, 0xD1, 0xFE, 0x9D, 0xD8, 0x51, 0xB3, 0xE7, 0xC6, 0x3B, 0xA3, 0x58,
    0xDD, 0xDE, 0x32, 0x9B, 0x98, 0x9A, 0x00, 0x49, 0xAB, 0x00, 0x6A, 0xAD, 0xD8, 0x0A, 0xAC, 0xE8,
    0xE3, 0xFF, 0xC2, 0x82, 0xD8, 0x94, 0xB5, 0x72, 0x5F, 0x2D, 0x72, 0xD5, 0xD9, 0x87, 0x43, 0xFC,
    0xF1, 0xA9, 0xC0, 0x2C, 0x60, 0xB2, 0xED, 0xBD, 0xEA, 0x7B, 0x03, 0x28, 0xD8, 0xD3, 0x65, 0x5E,
    0xD9, 0xB1, 0x82, 0xBE, 0x6C, 0x5B, 0x03, 0xB5, 0xC8, 0x4B, 0x75, 0x34, 0x40, 0x4D, 0x9A, 0x65,
    0xD6, 0xE6, 0x49, 0xDF, 0x5A, 0x28, 0xF5, 0x2A, 0xEF, 0x35, 0x3C, 0xA5, 0x4A, 0x45, 0x30, 0x14,
    0xFB, 0x37, 0xAE, 0x8F, 0x97, 0xC1, 0x92, 0x9B, 0x01, 0x2B, 0x16, 0xEA, 0x21, 0xA0, 0x1A, 0xDD,
    0xDF, 0xC4, 0xBA, 0x05, 0xBC, 0xC7, 0x4E, 0x8F, 0x9A, 0x50, 0xE4, 0x22, 0x58, 0x0D, 0xFB, 0xCB,
];

const RNG_BIT_LEN: usize = 1024;
const KEY_INFO: &str = "SRP Demo Key Information";


fn main() {

    let g = 2.to_bigint().unwrap();

    #[allow(non_snake_case)]
    let N = BigInt::from_bytes_be(Sign::Plus, &_N_BYTES);

    let reader = io::stdin();
    let mut out = io::stdout();

    println!("N (base64) : {}", base64::encode(&_N_BYTES)); 
    out.write("Enter Username: ".as_bytes()).unwrap();

    let mut input = String::new();
    input.clear(); out.flush().unwrap();
    reader.read_line(&mut input).unwrap();
    let username = input.trim().to_string();    

    out.write("Enter Password: ".as_bytes()).unwrap();
    input.clear(); out.flush().unwrap();
    reader.read_line(&mut input).unwrap();
    let password = input.trim().to_string();

    let identity = format!("{}:{}", username, password);

    println!("Identity : '{}'", identity);

    let mut rng = OsRng::new().unwrap();

    let a = rng.gen_biguint(RNG_BIT_LEN).to_bigint().unwrap();

    #[allow(non_snake_case)]    
    let A = g.pos_mod_pow(&a, &N);

    println!("a (hex) : {}\nA (hex) : {}\nA (base64) : {}", 
             a.to_str_radix(16),
             A.to_str_radix(16), 
             base64::encode(&A.to_bytes_be().1));

    out.write("Enter Salt: ".as_bytes()).unwrap();
    input.clear(); out.flush().unwrap();
    reader.read_line(&mut input).unwrap();
    let salt_str = input.trim().to_string();
    let salt_bytes = base64::decode(&salt_str).unwrap();

    let mut mac = HmacSha256::new_varkey(&salt_bytes).unwrap();
    mac.input(&identity.as_bytes());
    let x = BigInt::from_bytes_be(Sign::Plus, &mac.result().code()); 
    let v = g.pos_mod_pow(&x, &N);

    println!("v (hex) : {}\nv (base64): {}",
        v.to_str_radix(16),
        base64::encode(&v.to_bytes_be().1));

    out.write("B: ".as_bytes()).unwrap();
    input.clear(); out.flush().unwrap();
    reader.read_line(&mut input).unwrap();
    
    #[allow(non_snake_case)]
    let B_str = input.trim().to_string();
    #[allow(non_snake_case)]
    let B_bytes = base64::decode(&B_str).unwrap();
    #[allow(non_snake_case)]
    let B = BigInt::from_bytes_be(Sign::Plus, &B_bytes);

    let mut mac = HmacSha256::new_varkey(&A.to_bytes_be().1).unwrap();
    mac.input(&B_bytes);
    let u = BigInt::from_bytes_be(Sign::Plus, &mac.result().code());

    let mut mac = HmacSha256::new_varkey(&g.to_bytes_be().1).unwrap();
    mac.input(&_N_BYTES);
    let k = BigInt::from_bytes_be(Sign::Plus, &mac.result().code());

    #[allow(non_snake_case)]    
    let S_c = (&B - &k * &g.pos_mod_pow(&x, &N)).pos_mod_pow(&(&a + &u * &x), &N);
    let mut okm = [0u8;16];
    let hk = Hkdf::<Sha256>::extract(Some(&u.to_bytes_be().1[..]), &S_c.to_bytes_be().1[..]);
    hk.expand(&KEY_INFO.as_bytes(), &mut okm).unwrap();
    let k1 = BigInt::from_bytes_be(Sign::Plus, &okm[..]);

    let mut mac = HmacSha256::new_varkey(&okm[..]).unwrap();
    mac.input(&A.pos_mod_pow(&B, &N).to_bytes_be().1);

    #[allow(non_snake_case)]    
    let M_1_bytes = mac.result().code();

    #[allow(non_snake_case)]    
    let M_1 = BigInt::from_bytes_be(Sign::Plus, &M_1_bytes);


    println!("B (hex): {}\nu: (hex): {}\nS_c (hex): {}\nS_c (base64): {}\nk1 (hex): {}\nM1 (hex): {}\nM1 (base64): {}",
        B.to_str_radix(16),
        u.to_str_radix(16),
        S_c.to_str_radix(16),
        base64::encode(&S_c.to_bytes_be().1),
        k1.to_str_radix(16),
        M_1.to_str_radix(16),
        base64::encode(&M_1_bytes));

    out.write("M2: ".as_bytes()).unwrap();
    input.clear(); out.flush().unwrap();
    reader.read_line(&mut input).unwrap();


    #[allow(non_snake_case)]    
    let M_2_s_str = input.trim().to_string();
    #[allow(non_snake_case)]    
    let M_2_s_bytes = base64::decode(&M_2_s_str).unwrap();
    #[allow(non_snake_case)]
    let M_2_s = BigInt::from_bytes_be(Sign::Plus, &M_2_s_bytes);

    let mut mac = HmacSha256::new_varkey(&okm[..]).unwrap();
    mac.reset();
    mac.input(&A.pos_mod_pow(&M_1, &N).to_bytes_be().1);


    #[allow(non_snake_case)]    
    let M_2_c_bytes = mac.result().code();

    #[allow(non_snake_case)]    
    let M_2_c = BigInt::from_bytes_be(Sign::Plus, &M_2_c_bytes);

    if &M_2_s == &M_2_c {
        println!("Authenticated!");
    } else {
        println!("Unauthenticated!");
    }
}
