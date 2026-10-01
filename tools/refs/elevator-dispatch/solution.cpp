#include <cstdlib>
#include <set>
using namespace std;

class Elevator {
    set<int> stops;
    int at = 0, moved = 0;
    bool up = true;

    // The nearest stop above (or below) the lift, or -1.
    int seek(bool above) {
        auto it = stops.lower_bound(at);
        if (above) return it == stops.end() ? -1 : *it;
        return it == stops.begin() ? -1 : *prev(it);
    }

public:
    Elevator(int floors) {}

    void call(int floor) {
        if (floor != at) stops.insert(floor);
    }

    int next() {
        int stop = seek(up);
        if (stop < 0) {
            stop = seek(!up);
            if (stop < 0) return -1;
            up = !up;
        }
        moved += abs(stop - at);
        at = stop;
        stops.erase(stop);
        return stop;
    }

    int travelled() { return moved; }
};
