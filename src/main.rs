use cucumber::{World as _, given, then, when};
use std::time::Duration;
use tokio::time::sleep;

#[derive(Debug, Default, cucumber::World)]
struct World {
    user: Option<String>,
    capacity: usize,
}

#[given(expr = "{word} is hungry")]
async fn someone_is_hungry(w: &mut World, user: String) {
    sleep(Duration::from_secs(2)).await;

    w.user = Some(user);
}

#[when(regex = r"^(?:he|she|they) eats? (\d+) cucumbers?$")]
async fn eat_cucumbers(w: &mut World, count: usize) {
    sleep(Duration::from_secs(2)).await;

    w.capacity += count;

    assert!(w.capacity < 4, "{} exploded!", w.user.as_ref().unwrap());
}

#[then("she is full")]
async fn is_full(w: &mut World) {
    sleep(Duration::from_secs(2)).await;

    assert_eq!(w.capacity, 3, "{} isn't full!", w.user.as_ref().unwrap());
}

#[tokio::main]
async fn main() {
    World::run("tests/features/readme").await;
}
