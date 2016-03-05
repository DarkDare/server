using Go = import "../../common/capnp/go.capnp";

$Go.package("capnp");
$Go.import("goshawkdb.io/server/capnp");

@0xbbc717d787db5c5f;

struct Configuration {
  clusterId          @0: Text;
  version            @1: UInt32;
  hosts              @2: List(Text);
  f                  @3: UInt8;
  maxRMCount         @4: UInt8;
  asyncFlush         @5: Bool;
  rms                @6: List(UInt32);
  rmsRemoved         @7: List(UInt32);
  fingerprints       @8: List(Data);
  union {
    transitioningTo :group {
      configuration  @9: Configuration;
      allHosts       @10: List(Text);
      newRMIds       @11: List(UInt32);
      pendingInstall @12: List(UInt32);
      pending        @13: List(ConditionPair);
    }
    stable           @14: Void;
  }
}

struct ConditionPair {
  rmId      @0: UInt32;
  condition @1: Condition;
}

struct Condition {
  union {
    and       @0: Conjunction;
    or        @1: Disjunction;
    generator @2: Generator;
  }
}

struct Conjunction {
  left  @0: Condition;
  right @1: Condition;
}

struct Disjunction {
  left  @0: Condition;
  right @1: Condition;
}

struct Generator {
  rmId      @0: UInt32;
  permLen   @1: UInt8;
  start     @2: UInt8;
  union {
    lenSimple          @3: UInt8;
    lenAdjustIntersect @4: List(UInt32);
  }
  includes  @5: Bool;
}
