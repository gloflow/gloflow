// Small seedable random number generator (SplitMix64). Same seed gives the same images.

pub struct Rng {
    state: u64,
}

impl Rng {
    pub fn new(seed: u64) -> Rng {
        Rng { state: seed }
    }

    fn next_u64(&mut self) -> u64 {
        self.state = self.state.wrapping_add(0x9E37_79B9_7F4A_7C15);
        let mut z = self.state;
        z = (z ^ (z >> 30)).wrapping_mul(0xBF58_476D_1CE4_E5B9);
        z = (z ^ (z >> 27)).wrapping_mul(0x94D0_49BB_1331_11EB);
        z ^ (z >> 31)
    }

    /// Random float in [0, 1).
    pub fn float(&mut self) -> f64 {
        (self.next_u64() >> 11) as f64 / (1u64 << 53) as f64
    }

    /// Random float in [lo, hi).
    pub fn uniform(&mut self, lo: f64, hi: f64) -> f64 {
        lo + (hi - lo) * self.float()
    }

    /// Random integer in [lo, hi], both ends included.
    pub fn int(&mut self, lo: i32, hi: i32) -> i32 {
        let count = (hi - lo + 1) as u64;
        lo + (self.next_u64() % count) as i32
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn int_stays_in_range_and_hits_both_ends() {
        let mut rng = Rng::new(1);
        let mut seen_lo = false;
        let mut seen_hi = false;
        for _ in 0..1000 {
            let v = rng.int(3, 6);
            assert!((3..=6).contains(&v));
            seen_lo |= v == 3;
            seen_hi |= v == 6;
        }
        assert!(seen_lo && seen_hi);
    }

    #[test]
    fn same_seed_same_numbers() {
        let mut a = Rng::new(42);
        let mut b = Rng::new(42);
        for _ in 0..10 {
            assert_eq!(a.float(), b.float());
        }
    }
}
