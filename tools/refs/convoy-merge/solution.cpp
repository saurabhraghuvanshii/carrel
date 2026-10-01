#include <vector>
using namespace std;

vector<int> survivors(vector<int>& robots) {
    vector<int> left;
    for (int r : robots) {
        bool alive = true;
        while (alive && r < 0 && !left.empty() && left.back() > 0) {
            if (left.back() < -r) {
                left.pop_back();
            } else {
                if (left.back() == -r) left.pop_back();
                alive = false;
            }
        }
        if (alive) left.push_back(r);
    }
    return left;
}
