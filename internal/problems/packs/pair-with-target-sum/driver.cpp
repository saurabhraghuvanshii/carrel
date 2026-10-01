#include <algorithm>
#include <climits>
#include <cmath>
#include <cstdlib>
#include <exception>
#include <iostream>
#include <map>
#include <queue>
#include <set>
#include <string>
#include <unordered_map>
#include <unordered_set>
#include <vector>
using namespace std;

#include "solution.cpp"

// Reads T cases from stdin, calls the learner's code, prints one line per case.
// Anything the learner prints with cout goes to stderr.
int main() {
    streambuf* real = cout.rdbuf();
    cout.rdbuf(cerr.rdbuf());
    ostream out(real);
    int t;
    cin >> t;
    for (int c = 0; c < t; c++) {
        int n;
        cin >> n;
        vector<int> nums(n);
        for (auto& x : nums) cin >> x;
        int target;
        cin >> target;
        try {
            vector<int> r = pairWithTarget(nums, target);
            if (r.size() != 2) {
                out << "ERROR bad return value\n";
            } else {
                out << min(r[0], r[1]) << ' ' << max(r[0], r[1]) << '\n';
            }
        } catch (const exception& e) {
            out << "ERROR " << e.what() << '\n';
        }
        out.flush();
    }
    return 0;
}
