import style from "./index.less";
import TemporaryQuery from "@/pages/DataAnalysis/TemporaryQuery";
import RealTimeTrafficFlow from "@/pages/DataAnalysis/RealTimeBusinessFlow";
import DataAnalysisNav from "@/pages/DataAnalysis/components/Nav";
import DataSourceManage from "@/pages/DataAnalysis/DataSourceManage";
import { useIntl, useModel } from "umi";
import { BigDataNavEnum } from "@/pages/DataAnalysis/service/enums";
import { useMemo } from "react";
import OfflineManager from "@/pages/DataAnalysis/OfflineManager";
import ScreeningRow from "@/pages/DataAnalysis/components/ScreeningRow";
import { Empty } from "antd";
import ManageNodeModal from "@/pages/DataAnalysis/components/NodeManage/ManageNodeModal";
import ManageFolderModal from "@/pages/DataAnalysis/components/NodeManage/ManageFolderModal";

const DataAnalysis = () => {
  const { navKey, currentInstances } = useModel("dataAnalysis");
  const i18n = useIntl();

  const NavContent = useMemo(() => {
    if (!currentInstances) {
      return (
        <div className={style.defaultPage}>
          <Empty
            image={Empty.PRESENTED_IMAGE_SIMPLE}
            description={i18n.formatMessage({ id: "datasource.draw.selected" })}
          />
        </div>
      );
    }
    switch (navKey) {
      case BigDataNavEnum.TemporaryQuery:
        return <TemporaryQuery />;
      case BigDataNavEnum.RealTimeTrafficFlow:
        return <RealTimeTrafficFlow />;
      case BigDataNavEnum.OfflineManage:
        return <OfflineManager />;
      case BigDataNavEnum.DataSourceManage:
        return <DataSourceManage />;
      default:
        return <></>;
    }
  }, [navKey, currentInstances]);

  return (
    <div className={style.main}>
      <div className={style.contentBox}>
        <DataAnalysisNav />
        <div className={style.content}>{NavContent}</div>
      </div>
      <div className={style.positionBox}>
        <ScreeningRow />
      </div>
      <ManageNodeModal />
      <ManageFolderModal />
    </div>
  );
};

export default DataAnalysis;
export { BigDataNavEnum } from "@/pages/DataAnalysis/service/enums";
